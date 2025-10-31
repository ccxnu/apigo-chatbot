-- Clean up all pending registrations to remove ones with incorrect timezone
DELETE FROM cht_pending_registrations WHERE pnd_verified = FALSE;
