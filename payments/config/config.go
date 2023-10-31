package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/zeebo/errs"
)

var configErr = errs.Class("configuration error")

// Name of the configuration file.
const envFilePath = "config.env"

// Names of environment variables.
const (
	serverAddress  = "SERVER_ADDR"
	serverPort     = "SERVER_PORT"
	dbAddress      = "DB_ADDR"
	dbPort         = "DB_PORT"
	dbName         = "DB_NAME"
	dbUser         = "DB_USER"
	dbPassword     = "DB_PASS"
	dbSchemaPath   = "DB_SCHEMA_PATH"
	hashSalt       = "HASH_SALT"
	tokenSignature = "TOKEN_SIGNATURE"
	tokenTTL       = "TOKEN_TTL_HOURS"
)

// Configuration default settings.
const (
	defaultServerAddress  = "127.0.0.1"
	defaultServerPort     = "8080"
	defaultDBAddress      = "127.0.0.1"
	defaultDBPort         = "5432"
	defaultDBName         = "payments"
	defaultDBUser         = "payments"
	defaultDBPassword     = "lthgfhjk"
	defaultDBSchemaPath   = "./db/schema.sql"
	defaultHashSalt       = "HaShSaLt"
	defaultTokenSignature = "Some_Token_Signature"
	defaultTokenTTL       = "24"
)

func init() {
	err := godotenv.Load(envFilePath)
	if err != nil {
		log.Println(configErr.Wrap(err))
		setDefaults()
		if err = writeConfig(envFilePath); err != nil {
			log.Fatal(configErr.Wrap(err))
		}
	}
}

func ServerAddress() string {
	return os.Getenv(serverAddress)
}

func ServerPort() string {
	return ":" + os.Getenv(serverPort)
}

func DBAddress() string {
	return os.Getenv(dbAddress)
}

func DBPort() string {
	return os.Getenv(dbPort)
}

func DBName() string {
	return os.Getenv(dbName)
}

func DBUser() string {
	return os.Getenv(dbUser)
}

func DBPassword() string {
	return os.Getenv(dbPassword)
}

func DBSchemaPath() string {
	return os.Getenv(dbSchemaPath)
}

func HashSalt() string {
	return os.Getenv(hashSalt)
}

func TokenSignature() string {
	return os.Getenv(tokenSignature)
}

func TokenTTL() time.Duration {
	hours, err := strconv.Atoi(os.Getenv(tokenTTL))
	if err != nil {
		return time.Duration(0)
	}
	return time.Duration(hours) * time.Hour
}

// Sets the configuration to defaults.
func setDefaults() {
	log.Println("setting configuration defaults")
	os.Setenv(serverAddress, defaultServerAddress)
	os.Setenv(serverPort, defaultServerPort)
	os.Setenv(dbAddress, defaultDBAddress)
	os.Setenv(dbPort, defaultDBPort)
	os.Setenv(dbName, defaultDBName)
	os.Setenv(dbUser, defaultDBUser)
	os.Setenv(dbPassword, defaultDBPassword)
	os.Setenv(dbSchemaPath, defaultDBSchemaPath)
	os.Setenv(hashSalt, defaultHashSalt)
	os.Setenv(tokenSignature, defaultTokenSignature)
	os.Setenv(tokenTTL, defaultTokenTTL)
}

// Creates the configuration file.
func writeConfig(path string) error {
	log.Println("creating config.env file")
	wdErr := errs.Class("write defaults error")
	params := os.O_WRONLY | os.O_CREATE | os.O_APPEND
	file, err := os.OpenFile(path, params, 0600)
	if err != nil {
		return wdErr.Wrap(err)
	}

	// Write defaults to the created file.
	log.Println("saving defaults")
	_, err = fmt.Fprint(
		file,
		serverAddress+":"+ServerAddress()+"\n",
		serverPort+ServerPort()+"\n",
		dbAddress+":"+DBAddress()+"\n",
		dbPort+":"+DBPort()+"\n",
		dbName+":"+DBName()+"\n",
		dbUser+":"+DBUser()+"\n",
		dbPassword+":"+DBPassword()+"\n",
		dbSchemaPath+":"+DBSchemaPath()+"\n",
		hashSalt+":"+HashSalt()+"\n",
		tokenSignature+":"+TokenSignature()+"\n",
		tokenTTL+":"+fmt.Sprint(int64(TokenTTL()/time.Hour))+"\n",
	)
	if err != nil {
		return wdErr.Wrap(err)
	}
	if err = file.Close(); err != nil {
		return wdErr.Wrap(err)
	}
	return nil
}
