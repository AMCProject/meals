package database

type DbScript struct {
	File        string
	Script      string
	Description string
}

var scripts = []DbScript{
	{
		Script:      version,
		Description: "db version",
	},
	{
		Script:      users,
		Description: "users table",
	},
	{
		Script:      meals,
		Description: "meals table",
	},
	{
		Script:      calendars,
		Description: "calendars table",
	},
	{
		Script:      addNameToCalendars,
		Description: "add name column to calendar",
	},
}
var version = `
CREATE TABLE IF NOT EXISTS db_version (
	version int NOT NULL
);

INSERT OR IGNORE INTO db_version (version) VALUES (0);
`

var users = `
CREATE TABLE IF NOT EXISTS users (
	id		 text      PRIMARY KEY,
	name	 text	   NOT NULL,
	mail     text	   NOT NULL,
	password text      NOT NULL
);`

var meals = `
CREATE TABLE IF NOT EXISTS meals (
	id 		     text	 NOT NULL,
	user_id	     text 	 NOT NULL,
	name	 	 text	 NOT NULL,
	description  text,
	image		 text,
	kcal         integer NOT NULL,
	type		 text	 NOT NULL,
	ingredients	 text	 NOT NULL,
	seasons      text    NOT NULL,
	PRIMARY KEY (id,user_id)
);`

var calendars = `
CREATE TABLE IF NOT EXISTS calendar (
	user_id		text   NOT NULL,
	meal_id 	text   NOT NULL,
	date	    text   NOT NULL,
	PRIMARY KEY (user_id,date)
);`

var addNameToCalendars = `
ALTER TABLE calendar ADD name text NOT NULL;
`
