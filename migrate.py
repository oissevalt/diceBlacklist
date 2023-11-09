# migrate.py

import json
import sqlite3
import time


def fail(message):
    print(f'\033[91m{str.rjust("error", 10)}\033[0m {message}')


def success(message):
    print(f'\033[92m{str.rjust("success", 10)}\033[0m {message}')


conn: sqlite3.Connection | None = None
cursor: sqlite3.Cursor | None = None


def insert_data(__id, __reason, __initial, __latest):
    try:
        cursor.execute('''
        INSERT INTO blacklist (id, reason, initial, latest) VALUES (?, ?, ?, ?)
        ''', (__id, __reason, __initial, __latest))
        conn.commit()
    except sqlite3.Error as insert_error:
        fail(insert_error)


try:
    conn = sqlite3.connect("blacklist.sqlite.db")
    cursor = conn.cursor()

    cursor.execute('''
    CREATE TABLE IF NOT EXISTS blacklist(
        id TEXT PRIMARY KEY NOT NULL,
        reason TEXT,
        initial INTEGER NOT NULL,
        latest INTEGER
    )
    ''')
    conn.commit()

except sqlite3.Error as database_error:
    fail(database_error)

finally:
    if conn is not None and cursor is not None:
        success("Database successfully connected")
    else:
        exit(1)


total, counter = 0, 0
try:
    with open("blacklist.json", "r") as file:
        json_data = json.load(file)

    total = len(json_data)
    for item in json_data:
        if item.get("rank") != -30:
            continue

        item_id = item.get("ID")
        initial = item["times"][0] if item["times"] else time.time()
        latest = item["times"][-1] if item["times"] else time.time()
        reason = json.dumps(item.get("reasons"), ensure_ascii=False)

        insert_data(item_id, reason, initial, latest)
        success(f"Item with ID {item_id}")
        counter += 1

except (IOError, json.JSONDecodeError) as error:
    fail(f"Failed to read or parse JSON file: {error}")

except KeyError as key:
    fail(f"Key {key} was missing in JSON object")

finally:
    print(f"Migrated {counter} out of {total} items")
    conn.close()

