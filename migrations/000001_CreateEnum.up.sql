-- Create a gender enum type
CREATE TYPE gender_enum as ENUM ('Male','Female', 'Others', 'Unknown');
CREATE TYPE document_enum as ENUM ('Citizenship','Passport','Driving License','National ID','Pan Card');