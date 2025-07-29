-- +migrate Down

DELETE FROM SessionLog;
DELETE FROM UserRoleAssignment;
DELETE FROM UserProfile;
DELETE FROM Location;
DELETE FROM Business;
DELETE FROM Role;