package store

//Interface that every DB pkg must comply for MrRedis
type DB interface {

	//Perform the inital setup of the database/KV store by creating DB/Namespace etc that are important running MrRedis
	Setup(config string) error

	//Check if the database is setup already or not for Redis Framework
	IsSetup() bool

	//Optionally used if the db provides any auth mechanism perform that will handle DB apis like Connect/Login/Authorize etc.,
	Login() error

	//Set the value for the Key , if the key does not exisist create one (Will be an Insert if we RDBMS is introduced)
	Set(Key string, Value string) error

	//Does this key exist in the db already
	IsKey(Key string) (bool, error)

	//Update a particular Key with the value only if the key is valid already, optionally try to lock the key aswell (Update in RDBMS)
	Update(Key string, Value string, Lock bool) error

	//Get the value for a particular key (Will be a Select for RDBMS)
	Get(Key string) (string, error)

	//Delete a particular key from the store (Will be DEL for RDBMS)
	Del(Key string) error

	//Section
	//Section is a DIR in etcd
	//Section will be a namespace in Redis
	//Section will be a Table in RDBMS
	//Create Section
	CreateSection(Key string) error

	//Delete Section
	DeleteSection(Key string) error

	//List the complete secton if recursive then list the entirer directory tree, only the key not the value
	ListSection(Key string, Recursive bool) ([]string, error)

	//Completly wipe out the DB/KV store about all the information pertaining to MrRedis
	CleanSlate() error
}
