package client_test

// You MUST NOT change these default imports.  ANY additional imports may
// break the autograder and everyone will be sad.

import (
	// Some imports use an underscore to prevent the compiler from complaining
	// about unused imports.
	_ "encoding/hex"
	_ "errors"
	_ "strconv"
	_ "strings"
	"testing"

	// A "dot" import is used here so that the functions in the ginko and gomega
	// modules can be used without an identifier. For example, Describe() and
	// Expect() instead of ginko.Describe() and gomega.Expect().

	// custom

	// custom end

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cs161-staff/project2-starter-code/client"
	userlib "github.com/cs161-staff/project2-userlib"
)

func TestSetupAndExecution(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Tests")
}

// ================================================
// Global Variables (feel free to add more!)
// ================================================
const defaultPassword = "password"
const emptyString = ""
const contentOne = "Bitcoin is Nick's favorite "
const contentTwo = "digital "
const contentThree = "cryptocurrency!"

// ================================================
// Describe(...) blocks help you organize your tests
// into functional categories. They can be nested into
// a tree-like structure.
// ================================================

var _ = Describe("Client Tests", func() {

	// A few user declarations that may be used for testing. Remember to initialize these before you
	// attempt to use them!
	var alice *client.User
	var bob *client.User
	var charles *client.User
	// var doris *client.User
	// var eve *client.User
	// var frank *client.User
	// var grace *client.User
	// var horace *client.User
	// var ira *client.User

	// These declarations may be useful for multi-session testing.
	var alicePhone *client.User
	var aliceLaptop *client.User
	var aliceDesktop *client.User

	var err error

	// A bunch of filenames that may be useful.
	aliceFile := "aliceFile.txt"
	bobFile := "bobFile.txt"
	charlesFile := "charlesFile.txt"
	// dorisFile := "dorisFile.txt"
	// eveFile := "eveFile.txt"
	// frankFile := "frankFile.txt"
	// graceFile := "graceFile.txt"
	// horaceFile := "horaceFile.txt"
	// iraFile := "iraFile.txt"

	BeforeEach(func() {
		// This runs before each test within this Describe block (including nested tests).
		// Here, we reset the state of Datastore and Keystore so that tests do not interfere with each other.
		// We also initialize
		userlib.DatastoreClear()
		userlib.KeystoreClear()
	})

	Describe("Test InitUser", func() {

		Specify("Empty Username Error", func() {
			userlib.DebugMsg("Initializing empty usernmae user.")
			alice, err = client.InitUser("", defaultPassword)
			Expect(err).ToNot(BeNil())
			userlib.DebugMsg(err.Error())
		})

		Specify("2 users w/ same username", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initializing second user Alice.")
			alice, err = client.InitUser("alice", "betterPassword")
			Expect(err).ToNot(BeNil())
			userlib.DebugMsg(err.Error())
		})

		Specify("username case sensitive no error", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("Alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initializing user alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())
		})
		
		Specify("username case sensitive error", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initializing user alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).ToNot(BeNil())
		})

		Specify("User w/ empty password", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", "")
			Expect(err).To(BeNil())
		})

		Specify("2 users w/ same password", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initializing second user Bob.")
			alice, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())
		})

	})

	Describe("Test GetUser", func() {

		Specify("Get nonexisitng user", func() {
			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).ToNot(BeNil())
		})

		Specify("Incorrect password", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("alice", "wrongPassword")
			Expect(err).ToNot(BeNil())
		})

		Specify("Tampered w/ user", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Tamper w/ user Alice.")
			for key := range userlib.DatastoreGetMap() {
				userlib.DatastoreSet(key, []byte("tampered"))
			}

			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).ToNot(BeNil())
			userlib.DebugMsg(err.Error())
		})

	})

	Describe("Test StoreFile", func() {
		Specify("Store file no error", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file...")
			err = alice.StoreFile("alice1.txt", []byte("hello world!"))
			Expect(err).To(BeNil())
		})

		Specify("Store empty file ", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file...")
			err = alice.StoreFile("alice1.txt", []byte(""))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice loading file...")
			loadedContent, err := alice.LoadFile("alice1.txt")
			Expect(err).To(BeNil())
			Expect(loadedContent).To(Equal([]byte("")))
		})

		Specify("Filename shouldn't be globally unique", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file...")
			err = alice.StoreFile("alice1.txt", []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initializing user Bob.")
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob storing file...")
			err = bob.StoreFile("alice1.txt", []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice loading file...")
			loadedContent1, err := alice.LoadFile("alice1.txt")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob loading file...")
			loadedContent2, err := bob.LoadFile("alice1.txt")
			Expect(err).To(BeNil())

			Expect(loadedContent1).NotTo(Equal(loadedContent2))
		})

		Specify("Store file overwrite", func() {
			originalContent := []byte("Hello world")
			newContent := []byte("hey robert")

			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file...")
			err = alice.StoreFile("alice1.txt", originalContent)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice overwriting file...")
			err = alice.StoreFile("alice1.txt", newContent)
			Expect(err).To(BeNil())
			
			userlib.DebugMsg("Alice loading file...")
			loadedContent, err := alice.LoadFile("alice1.txt")
			Expect(err).To(BeNil())
			Expect(loadedContent).To(Equal(newContent))
		})
	})

	Describe("Test LoadFile", func() {
		Specify("File doesn't exist", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice loading file...")
			_, err := alice.LoadFile("alice1.txt")
			Expect(err).ToNot(BeNil())
		})


		Specify("Store/Load file w/ no error", func() {
			originalContent := []byte("hello world!")

			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file...")
			err = alice.StoreFile("alice1.txt", originalContent)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice loading file...")
			loadedContent, err := alice.LoadFile("alice1.txt")
			Expect(err).To(BeNil())
			Expect(loadedContent).To(Equal(originalContent))
		})

		Specify("Tamper with file. Expect err", func() {
			originalContent := []byte("hello world!")

			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			keys := getSeenKeys()

			userlib.DebugMsg("Alice storing file...")
			err = alice.StoreFile("alice1.txt", originalContent)
			Expect(err).To(BeNil())

			tamperUnseenKeys(keys)

			userlib.DebugMsg("Alice loading file...")
			_, err := alice.LoadFile("alice1.txt")
			Expect(err).NotTo(BeNil())
			userlib.DebugMsg(err.Error())
		})

		Specify("Multiple instances access same file", func() {
			originalContent := []byte("hello world!")

			userlib.DebugMsg("Initializing user Alice.")
			alicePhone, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user AliceLaptop.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())
			Expect(&alicePhone).To(Equal(&aliceLaptop))

			userlib.DebugMsg("AlicePhone storing file...")
			err = alicePhone.StoreFile("alice1.txt", originalContent)
			Expect(err).To(BeNil())
			Expect(&alicePhone).To(Equal(&aliceLaptop))

			userlib.DebugMsg("AliceLaptop loading file...")
			loadedContent, err := aliceLaptop.LoadFile("alice1.txt")
			Expect(err).To(BeNil())
			Expect(loadedContent).To(Equal(originalContent))
		})
	})

	Describe("AppendFile Tests", func() {
		Specify("Append File w/o err", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentTwo)
			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentThree)
			err = alice.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())
		})

		Specify("File doesn't exists", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentThree)
			err = alice.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).ToNot(BeNil())
		})

		Specify("Tamper with file before append", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			keys := getSeenKeys()

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			tamperUnseenKeys(keys)

			userlib.DebugMsg("Appending file data: %s", contentTwo)
			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())
		})

		Specify("Multiple instances appending to same file 1", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alicePhone, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user AliceLaptop.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("AlicePhone storing file...")
			err = alicePhone.StoreFile("alice1.txt", []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("AlicePhone appending to file...")
			err = alicePhone.AppendToFile("alice1.txt", []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("AliceLaptop loading file...")
			loadedContent, err := aliceLaptop.LoadFile("alice1.txt")
			Expect(err).To(BeNil())
			Expect(loadedContent).To(Equal([]byte(contentOne + contentTwo)))
		})

		Specify("Multiple instances appending to same file 2", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alicePhone, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user AliceLaptop.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("AlicePhone storing file...")
			err = alicePhone.StoreFile("alice1.txt", []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("AliceLaptop appending to file...")
			err = aliceLaptop.AppendToFile("alice1.txt", []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("AliceLaptop loading file...")
			loadedContent, err := aliceLaptop.LoadFile("alice1.txt")
			Expect(err).To(BeNil())
			Expect(loadedContent).To(Equal([]byte(contentOne + contentTwo)))

			userlib.DebugMsg("AliceLaptop loading file...")
			loadedContent, err = alicePhone.LoadFile("alice1.txt")
			Expect(err).To(BeNil())
			Expect(loadedContent).To(Equal([]byte(contentOne + contentTwo)))
		})

		Specify("Two different users appending to same file", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initializing user Bob.")
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file...")
			err = alice.StoreFile("alice1.txt", []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice invites bob")
			invitationPtr, err := alice.CreateInvitation("alice1.txt", "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepts invitation from Alice")
			err = bob.AcceptInvitation("alice", invitationPtr, "alice1.txt")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice appending to file...")
			err = alice.AppendToFile("alice1.txt", []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob loading file...")
			loadedContent, err := bob.LoadFile("alice1.txt")
			Expect(err).To(BeNil())
			Expect(loadedContent).To(Equal([]byte(contentOne + contentTwo)))
		})

		Specify("Two different users appending 100 times", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())
	
			userlib.DebugMsg("Alice storing file...")
			err = alice.StoreFile("alice1.txt", []byte(contentOne))
			Expect(err).To(BeNil())
	
			userlib.DebugMsg("Initializing user Bob.")
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())
	
			userlib.DebugMsg("Alice invites bob")
			invitationPtr, err := alice.CreateInvitation("alice1.txt", "bob")
			Expect(err).To(BeNil())
	
			userlib.DebugMsg("Bob accepts invitation from Alice")
			err = bob.AcceptInvitation("alice", invitationPtr, "alice1.txt")
			Expect(err).To(BeNil())
	
			i := 0
			for i < 100 {
				userlib.DebugMsg("Bob appending to file...")
				err = bob.AppendToFile("alice1.txt", []byte(contentOne))
				Expect(err).To(BeNil())
				i = i +1
			}

			userlib.DebugMsg("Alice storing file...")
			err = alice.StoreFile("alice1.txt", []byte(contentTwo))
			Expect(err).To(BeNil())
	
			userlib.DebugMsg("Bob loading file...")
			loadedContent, err := bob.LoadFile("alice1.txt")
			Expect(err).To(BeNil())
			Expect(loadedContent).To(Equal([]byte(contentTwo)))
		})

		Specify("Append File efficiency test", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			// Helper function to measure bandwidth of a particular operation
			measureBandwidth := func(probe func()) (bandwidth int) {
				before := userlib.DatastoreGetBandwidth()
				probe()
				after := userlib.DatastoreGetBandwidth()
				return after - before
			}

			base := measureBandwidth(func() {
				i := 0
				for i < 1 {
					userlib.DebugMsg("Appending file data: %s", "a")
					err = alice.AppendToFile(aliceFile, []byte("a"))
					Expect(err).To(BeNil())
					i = i +1
				}
			})

			bw := measureBandwidth(func() {
				i := 0
				for i < 100 {
					userlib.DebugMsg("Appending file data: %s", "a")
					err = alice.AppendToFile(aliceFile, []byte("a"))
					Expect(err).To(BeNil())
					i = i +1
				}
			})
			Expect(bw).To(Equal(base*100))
		})

	})

	Describe("Create/Accept Invitation Tests", func() {
		Specify("create and invitation w/o error", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initializing user Bob.")
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice invites bob")
			invitationPtr, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepts invitation from Alice")
			err = bob.AcceptInvitation("alice", invitationPtr, aliceFile)
			Expect(err).To(BeNil())
		})

		Specify("create invitation without creating file", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initializing user Bob.")
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice invites bob")
			_, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).NotTo(BeNil())
			userlib.DebugMsg(err.Error())
		})

		Specify("create invitation for inexistent user", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice invites bob")
			_, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).NotTo(BeNil())
			userlib.DebugMsg(err.Error())
		})

		Specify("create/accept invitation multiple users", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initializing user Bob.")
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initializing user Dave.")
			dave, err := client.InitUser("dave", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice invites bob")
			invitationPtr, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepts invitation from Alice")
			err = bob.AcceptInvitation("alice", invitationPtr, aliceFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob invites Dave")
			invitationPtr, err = bob.CreateInvitation(aliceFile, "dave")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Dave accepts invitation from Bob")
			err = dave.AcceptInvitation("bob", invitationPtr, aliceFile)
			Expect(err).To(BeNil())
		})

		Specify("create/accept invitation multiple users with tampering err", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initializing user Bob.")
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initializing user Dave.")
			dave, err := client.InitUser("dave", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			keys := getSeenKeys()

			userlib.DebugMsg("Alice invites bob")
			invitationPtr, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepts invitation from Alice")
			err = bob.AcceptInvitation("alice", invitationPtr, aliceFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob invites Dave")
			invitationPtr, err = bob.CreateInvitation(aliceFile, "dave")
			Expect(err).To(BeNil())

			tamperUnseenKeys(keys)

			userlib.DebugMsg("Dave accepts invitation from Bob")
			err = dave.AcceptInvitation("bob", invitationPtr, aliceFile)
			Expect(err).NotTo(BeNil())
			userlib.DebugMsg(err.Error())
		})

		Specify("accept invitation with already existing filename", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initializing user Bob.")
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob Storing file data: %s", contentOne)
			err = bob.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice invites bob")
			invitationPtr, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepts invitation from Alice")
			err = bob.AcceptInvitation("alice", invitationPtr, aliceFile)
			Expect(err).NotTo(BeNil())
			userlib.DebugMsg(err.Error())
		})

		Specify("accept invitation from non existent user", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initializing user Bob.")
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice invites bob")
			invitationPtr, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepts invitation from Alice")
			err = bob.AcceptInvitation("dave", invitationPtr, aliceFile)
			Expect(err).NotTo(BeNil())
			userlib.DebugMsg(err.Error())
		})

		Specify("create/accept invitation with multiple instances", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alicePhone, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initializing user Bob.")
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user AliceLaptop.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("AlicePhone storing file...")
			err = alicePhone.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("AliceLaptop invites bob")
			invitationPtr, err := aliceLaptop.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepts invitation from AliceLaptop")
			err = bob.AcceptInvitation("alice", invitationPtr, aliceFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("AlicePhone appending to file...")
			err = alicePhone.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob appending to file...")
			err = bob.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("AliceLaptop loading file...")
			loadedContent, err := aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(loadedContent).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Bob loading file...")
			loadedContent, err = bob.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(loadedContent).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("AlicePhone revokes invitation from bob")
			err = alicePhone.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())
			
			userlib.DebugMsg("Bob tries to append to file...")
			err = bob.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).NotTo(BeNil())
			userlib.DebugMsg(err.Error())
		})
	})

	Describe("Revoke Invitation Tests", func() {
		Specify("File doesn't exist", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initializing user Bob.")
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice revokes invitation from bob")
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).NotTo(BeNil())
			userlib.DebugMsg(err.Error())
		})

		Specify("filename is not currently shared with recipientUsername", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initializing user Bob.")
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice revokes invitation from bob")
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).NotTo(BeNil())
			userlib.DebugMsg(err.Error())
		})

		Specify("Successful revocation", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initializing user Bob.")
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice invites bob")
			invitationPtr, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepts invitation from Alice")
			err = bob.AcceptInvitation("alice", invitationPtr, aliceFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice revokes invitation from bob")
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())
		})

		Specify("Undefined behavior revoke", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initializing user Bob.")
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initializing user Dave.")
			dave, err := client.InitUser("dave", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice invites bob")
			invitationPtr, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepts invitation from Alice")
			err = bob.AcceptInvitation("alice", invitationPtr, aliceFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob invites dave")
			invitationPtr, err = bob.CreateInvitation(aliceFile, "dave")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Dave accepts invitation from Bob")
			err = dave.AcceptInvitation("bob", invitationPtr, aliceFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob revokes invitation from dave")
			err = bob.RevokeAccess(aliceFile, "dave")
			userlib.DebugMsg(err.Error())
			Expect(err).ToNot(BeNil())
			
		})

		Specify("Accept invitation after revocation", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initializing user Bob.")
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice invites bob")
			invitationPtr, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice revokes invitation from bob")
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepts invitation from Alice")
			err = bob.AcceptInvitation("alice", invitationPtr, aliceFile)
			Expect(err).NotTo(BeNil())
			userlib.DebugMsg(err.Error())
		})
	})

	Describe("Fuzz Testing", func() {
		Specify("Random Tampering attack", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initializing user Bob.")
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			tamperAllKeys()

			userlib.DebugMsg("Alice storing file...")
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).ToNot(BeNil())
			userlib.DebugMsg(err.Error())

			userlib.DebugMsg("Alice invites bob")
			invitationPtr, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).ToNot(BeNil())
			userlib.DebugMsg(err.Error())

			userlib.DebugMsg("Bob accepts invitation from Alice")
			err = bob.AcceptInvitation("alice", invitationPtr, aliceFile)
			Expect(err).ToNot(BeNil())
			userlib.DebugMsg(err.Error())

			userlib.DebugMsg("Alice appending to file...")
			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())
			userlib.DebugMsg(err.Error())

			userlib.DebugMsg("Bob appending to file...")
			err = bob.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).ToNot(BeNil())
			userlib.DebugMsg(err.Error())

			userlib.DebugMsg("Alice loading file...")
			_, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).ToNot(BeNil())
			userlib.DebugMsg(err.Error())

			userlib.DebugMsg("Bob loading file...")
			_, err = bob.LoadFile(aliceFile)
			Expect(err).ToNot(BeNil())
			userlib.DebugMsg(err.Error())

			userlib.DebugMsg("AlicePhone revokes invitation from bob")
			err = alicePhone.RevokeAccess(aliceFile, "bob")
			Expect(err).ToNot(BeNil())
			userlib.DebugMsg(err.Error())
			
			userlib.DebugMsg("Bob tries to append to file...")
			err = bob.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).NotTo(BeNil())
			userlib.DebugMsg(err.Error())
		})
	})
	
	Describe("Basic Tests", func() {

		Specify("Basic Test: Testing InitUser/GetUser on a single user.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())
		})

		Specify("Basic Test: Testing Single User Store/Load/Append.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentTwo)
			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentThree)
			err = alice.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing Create/Accept Invite Functionality with multiple users and multiple instances.", func() {
			userlib.DebugMsg("Initializing users Alice (aliceDesktop) and Bob.")
			aliceDesktop, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting second instance of Alice - aliceLaptop")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop storing file %s with content: %s", aliceFile, contentOne)
			err = aliceDesktop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceLaptop creating invite for Bob.")
			invite, err := aliceLaptop.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite from Alice under filename %s.", bobFile)
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob appending to file %s, content: %s", bobFile, contentTwo)
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop appending to file %s, content: %s", aliceFile, contentThree)
			err = aliceDesktop.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that aliceDesktop sees expected file data.")
			data, err := aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that aliceLaptop sees expected file data.")
			data, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that Bob sees expected file data.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Getting third instance of Alice - alicePhone.")
			alicePhone, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that alicePhone sees Alice's changes.")
			data, err = alicePhone.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing Revoke Functionality", func() {
			userlib.DebugMsg("Initializing users Alice, Bob, and Charlie.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))

			userlib.DebugMsg("Alice creating invite for Bob for file %s, and Bob accepting invite under name %s.", aliceFile, bobFile)

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob can load the file.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Bob creating invite for Charles for file %s, and Charlie accepting invite under name %s.", bobFile, charlesFile)
			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Charles can load the file.")
			data, err = charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Alice revoking Bob's access from %s.", aliceFile)
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob/Charles lost access to the file.")
			_, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())

			_, err = charles.LoadFile(charlesFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Checking that the revoked users cannot append to the file.")
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())

			err = charles.AppendToFile(charlesFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())
		})
	})
})

//Helper Functions for testing
func getSeenKeys() []userlib.UUID {
	keys := make([]userlib.UUID, len(userlib.DatastoreGetMap()))
	i := 0
	for key, _ := range userlib.DatastoreGetMap() {
		keys[i] = key
		i++
	}
	return keys
}

func tamperUnseenKeys(keys []userlib.UUID) {
	for key, _ := range userlib.DatastoreGetMap() {
		var result = false
		for _, x := range keys {
			if x == key {
				result = true
				break
			}
		}
		if !result {
			userlib.DatastoreSet(key, []byte("tampered"))
		}
	}
}

func tamperAllKeys() {
	for key, _ := range userlib.DatastoreGetMap() {
		userlib.DatastoreSet(key, []byte("tampered"))
	}
}
