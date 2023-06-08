package client

// CS 161 Project 2

// You MUST NOT change these default imports. ANY additional imports
// may break the autograder!

import (
	"encoding/json"

	userlib "github.com/cs161-staff/project2-userlib"
	"github.com/google/uuid"

	// hex.EncodeToString(...) is useful for converting []byte to string

	// Useful for string manipulation

	// Useful for formatting strings (e.g. `fmt.Sprintf`).
	"fmt"

	// Useful for creating new error messages to return using errors.New("...")
	"errors"

	// Optional.
	_ "strconv"
)

// This serves two purposes: it shows you a few useful primitives,
// and suppresses warnings for imports not being used. It can be
// safely deleted!
func someUsefulThings() {

	// Creates a random UUID.
	randomUUID := uuid.New()

	// Prints the UUID as a string. %v prints the value in a default format.
	// See https://pkg.go.dev/fmt#hdr-Printing for all Golang format string flags.
	userlib.DebugMsg("Random UUID: %v", randomUUID.String())

	// Creates a UUID deterministically, from a sequence of bytes.
	hash := userlib.Hash([]byte("user-structs/alice"))
	deterministicUUID, err := uuid.FromBytes(hash[:16])
	if err != nil {
		// Normally, we would `return err` here. But, since this function doesn't return anything,
		// we can just panic to terminate execution. ALWAYS, ALWAYS, ALWAYS check for errors! Your
		// code should have hundreds of "if err != nil { return err }" statements by the end of this
		// project. You probably want to avoid using panic statements in your own code.
		panic(errors.New("An error occurred while generating a UUID: " + err.Error()))
	}
	userlib.DebugMsg("Deterministic UUID: %v", deterministicUUID.String())

	// Declares a Course struct type, creates an instance of it, and marshals it into JSON.
	type Course struct {
		name      string
		professor []byte
	}

	course := Course{"CS 161", []byte("Nicholas Weaver")}
	courseBytes, err := json.Marshal(course)
	if err != nil {
		panic(err)
	}

	userlib.DebugMsg("Struct: %v", course)
	userlib.DebugMsg("JSON Data: %v", courseBytes)

	// Generate a random private/public keypair.
	// The "_" indicates that we don't check for the error case here.
	var pk userlib.PKEEncKey
	var sk userlib.PKEDecKey
	pk, sk, _ = userlib.PKEKeyGen()
	userlib.DebugMsg("PKE Key Pair: (%v, %v)", pk, sk)

	// Here's an example of how to use HBKDF to generate a new key from an input key.
	// Tip: generate a new key everywhere you possibly can! It's easier to generate new keys on the fly
	// instead of trying to think about all of the ways a key reuse attack could be performed. It's also easier to
	// store one key and derive multiple keys from that one key, rather than
	originalKey := userlib.RandomBytes(16)
	derivedKey, err := userlib.HashKDF(originalKey, []byte("mac-key"))
	if err != nil {
		panic(err)
	}
	userlib.DebugMsg("Original Key: %v", originalKey)
	userlib.DebugMsg("Derived Key: %v", derivedKey)

	// A couple of tips on converting between string and []byte:
	// To convert from string to []byte, use []byte("some-string-here")
	// To convert from []byte to string for debugging, use fmt.Sprintf("hello world: %s", some_byte_arr).
	// To convert from []byte to string for use in a hashmap, use hex.EncodeToString(some_byte_arr).
	// When frequently converting between []byte and string, just marshal and unmarshal the data.
	//
	// Read more: https://go.dev/blog/strings

	// Here's an example of string interpolation!
	_ = fmt.Sprintf("%s_%d", "file", 1)
}

// This is the type definition for the User struct.
// A Go struct is like a Python or Java class - it can have attributes
// (e.g. like the Username attribute) and methods (e.g. like the StoreFile method below).
type User struct {
	Username string
    Password string
	SecretKey userlib.PKEDecKey
    SignKey userlib.DSSignKey
    OwnedFileMap map[string]userlib.UUID  //Map of file name to file pointer only holds owned files
    SentSharedInvitation map[string]userlib.UUID //Map of file name+username to invitations we have sent
    AcceptedSharedInvitations map[string]userlib.UUID // Map of file name to invitiation only files not owned

	// You can add other attributes here if you want! But note that in order for attributes to
	// be included when this struct is serialized to/from JSON, they must be capitalized.
	// On the flipside, if you have an attribute that you want to be able to access from
	// this struct's methods, but you DON'T want that value to be included in the serialized value
	// of this struct that's stored in datastore, then you can use a "private" variable (e.g. one that
	// begins with a lowercase letter).
}

type AuthenticatedEncItem struct {
    Ciphertext []byte
    Tag []byte
}

type FileBlob struct {
    Content []byte
    NextContentPointer userlib.UUID
}

type FilePointer struct {
    FrontPointer userlib.UUID //pointer to first FileBlob
    EndPointer userlib.UUID	//pointer to last FileBlob
}

type Invitation struct {
    Owner string
    FilePointerUUID userlib.UUID
	Senders []string
}

func InitUser(username string, password string) (userdataptr *User, err error) {

	// err if empty username
	if username == "" { return nil, errors.New("username cannot be empty")}
	// err if username already exists
	var _, exists = userlib.KeystoreGet(username + "pk")
	if exists { return nil, errors.New("username already taken")}

	// Genererate public encryption keys
	var pk userlib.PKEEncKey
	var sk userlib.PKEDecKey
	pk, sk, _ = userlib.PKEKeyGen()

	// Generate digital signature
	var signK userlib.DSSignKey
	var verifyK userlib.DSVerifyKey
	signK, verifyK, _ = userlib.DSKeyGen()

	// Store to the keystore
    userlib.KeystoreSet(username + "pk", pk)
    userlib.KeystoreSet(username + "verify", verifyK)	

	// create a user
    var user User = User{username, password, sk, signK, make(map[string]userlib.UUID), make(map[string]userlib.UUID), make(map[string]userlib.UUID)}

	// generate and store salt. unique for each user
    salt := userlib.RandomBytes(16)
	saltHash := userlib.Hash([]byte(username + "salt"))
    saltUUID, err := uuid.FromBytes(saltHash[:16])
	if err != nil { return nil, err}
    userlib.DatastoreSet(saltUUID, salt)

	//Store user to datstore
	err = storeUserToDataStore(&user, salt)
	if err != nil { return nil, err}

	return &user, nil
}

func GetUser(username string, password string) (userdataptr *User, err error) {
	// Get salt from datastore
	salt, err := getSaltFromDataStore(username)
	if err != nil { return nil, err}

	// Generate keys to decrypt the user
	sourceKey := userlib.Argon2Key([]byte(password), salt, 16)
	userUUIDKey, err := userlib.HashKDF(sourceKey, []byte(username + "identifier"))
	if err != nil { return nil, err}
	userUUID, err := uuid.FromBytes(userUUIDKey[:16])
	if err != nil { return nil, err}

	// Get  encrypted user from dataStore
	authenticatedEncItemBytes, exists := userlib.DatastoreGet(userUUID)
	if (!exists) { return nil, errors.New("incorrect password")}
	var authenticatedEncItem AuthenticatedEncItem
	err = json.Unmarshal(authenticatedEncItemBytes, &authenticatedEncItem)
	if err != nil { return nil, err}

	// Check MAC of encrypted user
	macKey, err := userlib.HashKDF(sourceKey, []byte(username + "MacKey"))
	if err != nil { return nil, err}
	mac, err := userlib.HMACEval(macKey[:16], authenticatedEncItem.Ciphertext)
	if err != nil { return nil, err}
	isequal := userlib.HMACEqual(mac, authenticatedEncItem.Tag)
	if !isequal { return nil, errors.New("Encrypted User was tampered with")}

	// Decrypt the encrypted user
	userEncKey, err := userlib.HashKDF(sourceKey, []byte(username + "Symmetric EncKey"))
	if err != nil { return nil, err}
	userBytes := userlib.SymDec(userEncKey[:16], authenticatedEncItem.Ciphertext)
	var decryptedUser User
	err = json.Unmarshal(userBytes, &decryptedUser)
	if err != nil { return nil, err}

	return &decryptedUser, nil
}

func (userdata *User) StoreFile(filename string, content []byte) (err error) {
	// Get the user from datastore
	userdata, err = GetUser(userdata.Username, userdata.Password)
	if err != nil {	return err }

	if invitationUUID, ok := userdata.AcceptedSharedInvitations[filename]; ok {
		//This means filename is in AcceptedSharedInvitations

		//Get invitationn from datastore
		invite, err := getInvitationFromDataStore(invitationUUID)
		if err != nil {	return err }
		
		filePointerUUID := invite.FilePointerUUID

		// Get owner's salt from datastore
		salt, err := getSaltFromDataStore(invite.Owner)
		if err != nil { return err}

		// Create and Store New FileBlob
		nextRandomUUID := uuid.New()
		fileBlobUUID := uuid.New()
		err = createAndStoreNewFileBlob(salt, content, nextRandomUUID, fileBlobUUID)
		if err != nil { return err }
		
		// Get FilePointer from datastore
		filePointer, err := getFilePointerFromDataStore(filePointerUUID, salt)
		if err != nil { return err }

		filePointer.FrontPointer = fileBlobUUID
		filePointer.EndPointer = nextRandomUUID

		// Store FilePointer to the datastore
		err = storeFilePointerToDataStore(filePointerUUID, filePointer, salt)
		if err != nil { return err }

		return nil
	} else { // filename is owned by the user
		// Get user's salt from datastore
		salt, err := getSaltFromDataStore(userdata.Username)
		if err != nil { return err}

		// Create and Store New FileBlob
		nextRandomUUID := uuid.New()
		fileBlobUUID := uuid.New()
		err = createAndStoreNewFileBlob(salt, content, nextRandomUUID, fileBlobUUID)
		if err != nil { return err }

		var filePointer FilePointer
		if filePointerUUID, ok := userdata.OwnedFileMap[filename]; ok {
			//This means filename is in ownedFileMap
			
			// Get FilePointer from datastore
			filePointer, err := getFilePointerFromDataStore(filePointerUUID, salt)
			if err != nil { return err }

			filePointer.FrontPointer = fileBlobUUID
			filePointer.EndPointer = nextRandomUUID

			// Store FilePointer to the datastore
			filePointerUUID = userdata.OwnedFileMap[filename]
			err = storeFilePointerToDataStore(filePointerUUID, filePointer, salt)
			if err != nil { return err }
		} else {
			// File doesn't exists; make filePointer; filename is not in ownedFileMap
			filePointer = FilePointer{fileBlobUUID, nextRandomUUID}
			
			//Derve filePointer UUID
			filePointerUUID = uuid.New() 
			userdata.OwnedFileMap[filename] = filePointerUUID

			// Store FilePointer to the datastore
			filePointerUUID = userdata.OwnedFileMap[filename]
			err = storeFilePointerToDataStore(filePointerUUID, filePointer, salt)
			if err != nil { return err }
		}
	
		// Store the user
		err = storeUserToDataStore(userdata, salt)
		if err != nil { return err}

		return nil
	}
}

func (userdata *User) AppendToFile(filename string, content []byte) error {
	// Get user from data
	userdata, err := GetUser(userdata.Username, userdata.Password)
	if err != nil {	return err }

	if invitationUUID, ok := userdata.AcceptedSharedInvitations[filename]; ok {
		//This means filename is in AcceptedSharedInvitations

		//Get invitationn from datastore
		invite, err := getInvitationFromDataStore(invitationUUID)
		if err != nil { return err}
		
		filePointerUUID := invite.FilePointerUUID

		// Get salt from datastore
		salt, err := getSaltFromDataStore(invite.Owner)
		if err != nil { return err}
			
		// Get FilePointer from datastore
		filePointer, err := getFilePointerFromDataStore(filePointerUUID, salt)
		if err != nil { return err }

		// Create and Store New FileBlob
		nextRandomUUID := uuid.New()
		fileBlobUUID := filePointer.EndPointer
		err = createAndStoreNewFileBlob(salt, content, nextRandomUUID, fileBlobUUID)
		if err != nil { return err }

		// change the endpointer
		filePointer.EndPointer = nextRandomUUID

		// Store FilePointer to the datastore
		err = storeFilePointerToDataStore(filePointerUUID, filePointer, salt)
		if err != nil { return err }
		return nil
	}
	if filePointerUUID, ok := userdata.OwnedFileMap[filename]; ok {
		//This means filename is in OwnedFileMap

		// Get salt from datastore
		salt, err := getSaltFromDataStore(userdata.Username)
		if err != nil { return err}
			
		// Get FilePointer from datastore
		filePointer, err := getFilePointerFromDataStore(filePointerUUID, salt)
		if err != nil { return err }

		// Create and Store New FileBlob
		nextRandomUUID := uuid.New()
		fileBlobUUID := filePointer.EndPointer
		err = createAndStoreNewFileBlob(salt, content, nextRandomUUID, fileBlobUUID)
		if err != nil { return err }

		// change the endpointer
		filePointer.EndPointer = nextRandomUUID

		// Store FilePointer to the datastore
		filePointerUUID = userdata.OwnedFileMap[filename]
		err = storeFilePointerToDataStore(filePointerUUID, filePointer, salt)
		if err != nil { return err }
		return nil
	}

	return errors.New("User does not have access to a file called " + filename)
}

func (userdata *User) LoadFile(filename string) (content []byte, err error) {
	// Get the user from datastore
	userdata, err = GetUser(userdata.Username, userdata.Password)
	if err != nil {	return nil, err }

	if invitationUUID, ok := userdata.AcceptedSharedInvitations[filename]; ok {
		//This means filename is in AcceptedSharedInvitations

		//Get invitationn from datastore
		invite, err := getInvitationFromDataStore(invitationUUID)
		if err != nil { return nil, err}
		
		filePointerUUID := invite.FilePointerUUID

		// Get salt from datastore
		salt, err := getSaltFromDataStore(invite.Owner)
		if err != nil { return nil, err}
			
		// Get FilePointer from datastore
		filePointer, err := getFilePointerFromDataStore(filePointerUUID, salt)
		if err != nil { return nil, err }

		// Recreate content from linked list of FileBlobs 
		content, err := constructContentFromFilePointer(filePointer, salt)
		if err != nil { return nil, err }
		return content, nil

	}

	if filePointerUUID, ok := userdata.OwnedFileMap[filename]; ok {
		//This means filename is in OwnedFileMap

		// Get salt from datastore
		salt, err := getSaltFromDataStore(userdata.Username)
		if err != nil { return nil, err}
			
		// Get FilePointer from datastore
		filePointer, err := getFilePointerFromDataStore(filePointerUUID, salt)
		if err != nil { return nil, err }

		// Recreate content from linked list of FileBlobs 
		content, err := constructContentFromFilePointer(filePointer, salt)
		if err != nil { return nil, err }
		return content, nil
	}
	return nil, errors.New("User does not have access to a file called " + filename)
}

func (userdata *User) CreateInvitation(filename string, recipientUsername string) (
	invitationPtr uuid.UUID, err error) {

	// check if recipient exists
	_, err = getSaltFromDataStore(recipientUsername)
	if err != nil {	return uuid.Nil, errors.New("Recipientt does not exist") }

	// Get the user from datastore
	userdata, err = GetUser(userdata.Username, userdata.Password)
	if err != nil {	return uuid.Nil, err }

	// Get salt from datastore
	salt, err := getSaltFromDataStore(userdata.Username)
	if err != nil { return uuid.Nil, err}

	if invitationUUID, ok := userdata.AcceptedSharedInvitations[filename]; ok {
		//This means filename is in AcceptedSharedInvitations
		
		// Get encrypted invitation
		_, exists := userlib.DatastoreGet(invitationUUID)

		if !exists { 
			delete(userdata.AcceptedSharedInvitations, filename)

			//Store user to datstore
			err = storeUserToDataStore(userdata, salt)
			if err != nil { return uuid.Nil, err}

			return uuid.Nil, errors.New("no longer have acccess to file")
		}
	
		//Get invitationn from datastore
		invite, err := getInvitationFromDataStore(invitationUUID)
		if err != nil { return uuid.Nil, err}

		//Add username to invitation sender list
		invite.Senders = append(invite.Senders, userdata.Username)

		// Store invitation to datastore
		err = storeInvitationToDataStore(invitationUUID, invite)
		if err != nil { return uuid.Nil, err}
		return invitationUUID, nil
	}

	if filePointerUUID, ok := userdata.OwnedFileMap[filename]; ok {

		var senders []string
		senders = append(senders, userdata.Username)
		invite := Invitation{userdata.Username, filePointerUUID, senders}

		invitationUUID := uuid.New()

		//Derive keys for invitation encryption
		sourceKey := userlib.Argon2Key([]byte(invitationUUID.String()), []byte(invitationUUID.String()), 16)
		invitationEncKey, err := userlib.HashKDF(sourceKey, []byte("invitation symmetric key"))
		if err != nil { return uuid.Nil, err}
		invitationMacKey, err := userlib.HashKDF(sourceKey, []byte("invitation mac key"))
		if err != nil { return uuid.Nil, err}

		//Encrypt invitation
		iv := userlib.RandomBytes(16)
		invitationBytes, err := json.Marshal(invite)
		if err != nil { return uuid.Nil, err}
		ciphertext := userlib.SymEnc(invitationEncKey[:16], iv, invitationBytes) // don't have to worry about collison

		mac, err := userlib.HMACEval(invitationMacKey[:16], ciphertext)
		if err != nil { return uuid.Nil, err}
		authenticatedEncItem := AuthenticatedEncItem{ciphertext, mac}
		
		//Store invitation
		authenticatedEncItemBytes, err := json.Marshal(authenticatedEncItem)
		if err != nil { return uuid.Nil, err}
		userlib.DatastoreSet(invitationUUID, authenticatedEncItemBytes)

		//Update sentSharedInvitation map
		userdata.SentSharedInvitation[filename+recipientUsername] = invitationUUID

		//Store user to datstore
		err = storeUserToDataStore(userdata, salt)
		if err != nil { return uuid.Nil, err}
		return invitationUUID, nil
	}

	return uuid.Nil, errors.New("User has no file named " + filename)
}

func (userdata *User) AcceptInvitation(senderUsername string, invitationPtr uuid.UUID, filename string) error {
	
	_, exists := userlib.DatastoreGet(invitationPtr)
	if !exists { 
		return errors.New("no longer have acccess to file")
	}

	// Get the user from datastore
	userdata, err := GetUser(userdata.Username, userdata.Password)
	if err != nil {	return err }

	if _, ok := userdata.AcceptedSharedInvitations[filename]; ok {
		//This means filename is in AcceptedSharedInvitations
		return errors.New("Already have a file named " + filename)
	}

	if _, ok := userdata.OwnedFileMap[filename]; ok {
		//This means filename is inOwnedFileMap
		return errors.New("Already have a file named " + filename)
	}

	invite, err := getInvitationFromDataStore(invitationPtr)
	if err != nil { return err}

	result := false
	for _, sender := range invite.Senders {
        if sender == senderUsername {
            result = true
            break
        }
    }

	if !result {
		return errors.New("Unable to verify invitation sender")
	}

	userdata.AcceptedSharedInvitations[filename] = invitationPtr

	// Get salt from datastore
	salt, err := getSaltFromDataStore(userdata.Username)
	if err != nil { return err}

	//Store user to datstore
	err = storeUserToDataStore(userdata, salt)
	if err != nil { return err}
	return nil
}

func (userdata *User) RevokeAccess(filename string, recipientUsername string) error {

	// Get the user from datastore
	userdata, err := GetUser(userdata.Username, userdata.Password)
	if err != nil {	return err }

	if _, ok := userdata.AcceptedSharedInvitations[filename]; ok {
		//This means filename is in AcceptedSharedInvitations
		//Do not have to worry about this way
		return errors.New("Untested behavior: only owner of file can revoke")
	}

	if _, ok := userdata.OwnedFileMap[filename]; ok {
		//This means filename is inOwnedFileMap
		if _,ok := userdata.SentSharedInvitation[filename+recipientUsername]; !ok {
			return errors.New("filename is not currently shared with recipientUsername")
		}

		invitationUUID := userdata.SentSharedInvitation[filename+recipientUsername]

		userlib.DatastoreDelete(invitationUUID)

		delete(userdata.SentSharedInvitation, filename+recipientUsername)

		// Get salt from datastore
		salt, err := getSaltFromDataStore(userdata.Username)
		if err != nil { return err}

		//Store user to datstore
		err = storeUserToDataStore(userdata, salt)
		if err != nil { return err}
		return nil
	}
	return errors.New("User does not have access to a file called: " + filename)
}

// Helper functions
func storeUserToDataStore(userdata *User, salt []byte) (err error) {
	// Generate keys to encrypt the user
    sourceKey := userlib.Argon2Key([]byte(userdata.Password), salt, 16)
	userEncKey, err := userlib.HashKDF(sourceKey, []byte(userdata.Username + "Symmetric EncKey"))
	if err != nil { return err}
	macKey, err := userlib.HashKDF(sourceKey, []byte(userdata.Username + "MacKey"))
	if err != nil { return err}
	iv := userlib.RandomBytes(16)

	// Encrypt user
    userBytes, err := json.Marshal(userdata)
	if err != nil { return err}
    ciphertext := userlib.SymEnc(userEncKey[:16], iv, userBytes) // userEncKey[:16] : don't have to worry about collison
    mac, err := userlib.HMACEval(macKey[:16], ciphertext)
	if err != nil { return err}
    authenticatedEncItem := AuthenticatedEncItem{ciphertext, mac}

	// Store encrypted user to the datastore
	userUUIDKey, err := userlib.HashKDF(sourceKey, []byte(userdata.Username + "identifier"))
	if err != nil { return err}
	userUUID, err := uuid.FromBytes(userUUIDKey[:16])
	if err != nil { return err}
	authenticatedEncItemBytes, err := json.Marshal(authenticatedEncItem)
	if err != nil { return err}
	userlib.DatastoreSet(userUUID, authenticatedEncItemBytes)

	return nil
}

func createAndStoreNewFileBlob(salt []byte, content []byte,
	 nextRandomUUID userlib.UUID, fileBlobUUID userlib.UUID) (err error) {

	fileBlob := FileBlob{content, nextRandomUUID}
	
	//Derive Keys for fileBlob encryption
	sourceKey := userlib.Argon2Key([]byte(fileBlobUUID.String()), salt, 16)
	fileEncKey, err := userlib.HashKDF(sourceKey, []byte("file symmetric key"))
	if err != nil { return err}
	fileMacKey, err := userlib.HashKDF(sourceKey, []byte("file mac key"))
	if err != nil { return err}
		
	//Encrypt fileBlob
	iv := userlib.RandomBytes(16)
	fileBlobBytes, err := json.Marshal(fileBlob)
	if err != nil { return err}
	ciphertext := userlib.SymEnc(fileEncKey[:16], iv, fileBlobBytes) // don't have to worry about collison
		
	mac, err := userlib.HMACEval(fileMacKey[:16], ciphertext)
	if err != nil { return err}
	authenticatedEncItem := AuthenticatedEncItem{ciphertext, mac}
	
	//Store fileBlob
	authenticatedEncItemBytes, err := json.Marshal(authenticatedEncItem)
	if err != nil { return err}
	userlib.DatastoreSet(fileBlobUUID, authenticatedEncItemBytes)

	return nil
}

func storeFilePointerToDataStore(filePointerUUID userlib.UUID, filePointer FilePointer, salt []byte) (err error) {
	//Derive keys for filePointer decryption
	sourceKey := userlib.Argon2Key([]byte(filePointerUUID.String()), salt, 16)
	filePointerEncKey, err := userlib.HashKDF(sourceKey, []byte("filePointer symmetric key"))
	if err != nil { return err}
	filePointerMacKey, err := userlib.HashKDF(sourceKey, []byte("filePointer mac key"))
	if err != nil { return err}

	//Encrypt filePointer
	iv := userlib.RandomBytes(16)
	filePointerBytes, err := json.Marshal(filePointer)
	if err != nil { return err}
	ciphertext := userlib.SymEnc(filePointerEncKey[:16], iv, filePointerBytes) // don't have to worry about collison

	mac, err := userlib.HMACEval(filePointerMacKey[:16], ciphertext)
	if err != nil { return err}
	authenticatedEncItem := AuthenticatedEncItem{ciphertext, mac}
	
	//Store filePointer
	authenticatedEncItemBytes, err := json.Marshal(authenticatedEncItem)
	if err != nil { return err}

	userlib.DatastoreSet(filePointerUUID, authenticatedEncItemBytes)

	return nil
}

func getSaltFromDataStore(username string) (salt []byte, err error) {
	saltHash := userlib.Hash([]byte(username + "salt"))
	saltUUID, err := uuid.FromBytes(saltHash[:16])
	if err != nil { return nil, err}
	salt, exists := userlib.DatastoreGet(saltUUID)
	if !exists { return nil, errors.New("user doesn't exists")}
	return salt, nil
}

func getFilePointerFromDataStore(filePointerUUID userlib.UUID, salt []byte) (filePointer FilePointer, err error) {
	//Get current encrypted filePointer
	authenticatedEncItemBytes, exists := userlib.DatastoreGet(filePointerUUID)
	if !exists { return FilePointer{}, errors.New("file doesn't exists or file has been tampered with")}
	var authenticatedEncItem AuthenticatedEncItem
	err = json.Unmarshal(authenticatedEncItemBytes, &authenticatedEncItem)
	if err != nil { return FilePointer{}, err}
	
	//Derive keys for filePointer decryption
	sourceKey := userlib.Argon2Key([]byte(filePointerUUID.String()), salt, 16)
	filePointerEncKey, err := userlib.HashKDF(sourceKey, []byte("filePointer symmetric key"))
	if err != nil { return FilePointer{}, err}
	filePointerMacKey, err := userlib.HashKDF(sourceKey, []byte("filePointer mac key"))
	if err != nil { return FilePointer{}, err}

	// Check mac
	mac, err := userlib.HMACEval(filePointerMacKey[:16], authenticatedEncItem.Ciphertext)
	if err != nil { return FilePointer{}, err}
	isequal := userlib.HMACEqual(mac, authenticatedEncItem.Tag)
	if !isequal { return FilePointer{}, errors.New("Encrypted file was tampered with")}

	// Decrypt the encrypted filePointer
	filePointerBytes := userlib.SymDec(filePointerEncKey[:16], authenticatedEncItem.Ciphertext)
	err = json.Unmarshal(filePointerBytes, &filePointer)
	if err != nil { return FilePointer{}, err}

	return filePointer, nil
}

func getInvitationFromDataStore(invitationUUID userlib.UUID) (invitePtr Invitation, err error) {
	//Get invitationn from datastore
	authenticatedEncItemBytes, exists := userlib.DatastoreGet(invitationUUID)
	if !exists { return Invitation{}, errors.New("Invitation doesn't exists")}

	var authenticatedEncItem AuthenticatedEncItem 
	err = json.Unmarshal(authenticatedEncItemBytes, &authenticatedEncItem)
	if err != nil { return Invitation{}, err}

	//Derive keys for invitation decryption
	sourceKey := userlib.Argon2Key([]byte(invitationUUID.String()), []byte(invitationUUID.String()), 16)
	invitationEncKey, err := userlib.HashKDF(sourceKey, []byte("invitation symmetric key"))
	if err != nil { return Invitation{}, err}
	invitationMacKey, err := userlib.HashKDF(sourceKey, []byte("invitation mac key"))
	if err != nil { return Invitation{}, err}

	// Check mac
	mac, err := userlib.HMACEval(invitationMacKey[:16], authenticatedEncItem.Ciphertext)
	if err != nil { return Invitation{}, err}
	isequal := userlib.HMACEqual(mac, authenticatedEncItem.Tag)
	if !isequal { return Invitation{}, errors.New("Encrypted invitation was tampered with")}

	// Decrypt the encrypted invitation
	invitationBytes := userlib.SymDec(invitationEncKey[:16], authenticatedEncItem.Ciphertext)
	var invite Invitation
	err = json.Unmarshal(invitationBytes, &invite)
	if err != nil { return Invitation{}, err}

	return invite, nil
}

func storeInvitationToDataStore(invitationUUID userlib.UUID, invite Invitation) (err error) {
	//Derive keys for invitation decryption
	sourceKey := userlib.Argon2Key([]byte(invitationUUID.String()), []byte(invitationUUID.String()), 16)
	invitationEncKey, err := userlib.HashKDF(sourceKey, []byte("invitation symmetric key"))
	if err != nil { return err}
	invitationMacKey, err := userlib.HashKDF(sourceKey, []byte("invitation mac key"))
	if err != nil { return err}

	//Encrypt invitation
	iv := userlib.RandomBytes(16)
	invitationBytes, err := json.Marshal(invite)
	if err != nil { return err}
	ciphertext := userlib.SymEnc(invitationEncKey[:16], iv, invitationBytes) // don't have to worry about collison

	mac, err := userlib.HMACEval(invitationMacKey[:16], ciphertext)
	if err != nil { return err}
	authenticatedEncItem := AuthenticatedEncItem{ciphertext, mac}
	
	//Store invitation
	authenticatedEncItemBytes, err := json.Marshal(authenticatedEncItem)
	if err != nil { return err}
	userlib.DatastoreSet(invitationUUID, authenticatedEncItemBytes)

	return nil
}

func constructContentFromFilePointer(filePointer FilePointer, salt []byte) (content []byte, err error) {
	content = make([]byte, 0)
	currentBlobUUID := filePointer.FrontPointer
	for currentBlobUUID != filePointer.EndPointer {
		
		// Get encrypted fileBlob
		authenticatedEncItemBytes, exists := userlib.DatastoreGet(currentBlobUUID)

		if !exists { return nil, errors.New("file blob has been tampered with")}
		var authenticatedEncItem AuthenticatedEncItem 
		err = json.Unmarshal(authenticatedEncItemBytes, &authenticatedEncItem)
		if err != nil { return nil, err}

		//Derive Key for fileBlob decryption
		sourceKey := userlib.Argon2Key([]byte(currentBlobUUID.String()), salt, 16)
		fileEncKey, err := userlib.HashKDF(sourceKey, []byte("file symmetric key"))
		if err != nil { return nil, err}
		fileMacKey, err := userlib.HashKDF(sourceKey, []byte("file mac key"))
		if err != nil { return nil, err}

		// Check mac
		mac, err := userlib.HMACEval(fileMacKey[:16], authenticatedEncItem.Ciphertext)
		if err != nil { return nil, err}
		isequal := userlib.HMACEqual(mac, authenticatedEncItem.Tag)
		if !isequal { return nil, errors.New("Encrypted file was tampered with")}

		// Decrypt the encrypted blob
		fileBlobBytes := userlib.SymDec(fileEncKey[:16], authenticatedEncItem.Ciphertext)
		var fileBlob FileBlob
		err = json.Unmarshal(fileBlobBytes, &fileBlob)
		if err != nil { return nil, err}

		newContent := fileBlob.Content
		content = append(content, newContent...)
		currentBlobUUID = fileBlob.NextContentPointer
	}
	return content, nil
}