Create Database db_latihan_pbp;

/

Create Table users(
	ID INT(20) NOT NULL AUTO_INCREMENT,
	Name VARCHAR(60) NOT NULL,
	Age INT(40) NOT NULL,
	Address VARCHAR(200) NOT NULL,
	Password VARCHAR(255) NOT NULL,
	PRIMARY KEY(ID)
)
