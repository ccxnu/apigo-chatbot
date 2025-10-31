#!/bin/bash
export PGPASSWORD='lo0G4Rfaw7gtHw0wvpm4aqi4'
psql -h localhost -p 5432 -U postgres -d chatbot_db -f /home/xrnon/Dev/ISTS/new-chatbot/apigo-chatbot/scripts/fix_dirty_migration.sql
