# This script clears all active Agents from Endar server database.
# Must be executed within the container itself.
# Please use with caution as there is no confirmation.
# Written by Queball1999 on 01/09/2025


import psycopg2
import os

SQLALCHEMY_DATABASE_URI = os.environ.get("SQLALCHEMY_DATABASE_URI")

def clear_agents():
    try:
        connection = psycopg2.connect(SQLALCHEMY_DATABASE_URI)
        cursor = connection.cursor()
        cursor.execute("DELETE FROM agents")
        connection.commit()
        cursor.close()
        connection.close()
        return True
    except Exception as e:
        print(f"[ERROR] {e}".strip())
        return False

print(f"[INFO] Connecting to the database server: {SQLALCHEMY_DATABASE_URI}")
if clear_agents():
    print(f"[INFO] Successfully cleared all agents from the database")
else:
    print(f"[ERROR] Unable to clear agents from the database")
    exit(1)
exit(0)