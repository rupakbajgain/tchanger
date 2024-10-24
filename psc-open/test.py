import requests
import sqlite3
import json
from pathlib import Path

p=Path(__file__).parents[0]/'main.db'
#print(p)
conn = sqlite3.connect(p)

def notify(title, body):
    try:
        requests.post('http://localhost:8080/notify',json={'title':title,'body':body})
    except:
        pass

def fetch_content(o):
    path='/gateway/psc/dashboard/fetch-basics'
    if o['ID']==2 or o['ID']==4 or o['ID']==5:
        path='/gateway/psc/getLiveAdvtData'
    headers={
        "Authorization": "Bearer "+o['auth'],
    }
    x=requests.get(o['baseUrl']+path, headers=headers)
    print(o['baseUrl']+path, x)
    if x.status_code!=401:
        o['json']=x.text
    return x.status_code

def login(o):
    headers = {
        "Content-Type": "multipart/form-data; boundary=---------------------------250594647310496953021971810313"
    }
    data=f"-----------------------------250594647310496953021971810313\r\nContent-Disposition: form-data; name=\"username\"\r\n\r\n{o['user']}\r\n-----------------------------250594647310496953021971810313\r\nContent-Disposition: form-data; name=\"password\"\r\n\r\n{o['password']}\r\n-----------------------------250594647310496953021971810313\r\nContent-Disposition: form-data; name=\"grant_type\"\r\n\r\npassword\r\n-----------------------------250594647310496953021971810313--\r\n"
    #print(data)
    path='/gateway/login'
    #print(o['baseUrl']+path)
    x=requests.post(o['baseUrl']+path, data, headers=headers)
    print(x)
    json=x.json()
    #print(json)
    access_token=json['access_token']
    o['auth']=access_token
    #save to db
    cursor = conn.cursor()
    sql=''' UPDATE urls
            SET authorization = ?
            WHERE ID= ?'''
    cursor.execute(sql, (o['auth'], o['ID']))
    conn.commit()

def checkSingleUpdate(o):
    if not o['auth']:
        login(o)
    if fetch_content(o)==401:
        login(o)
        if fetch_content(o)==401:
            return
    if o['lastResponse']==o['json']:
        return
    #save to db
    cursor = conn.cursor()
    sql=''' UPDATE urls
            SET lastResponse = ?
            WHERE ID= ?'''
    cursor.execute(sql, (o['json'], o['ID']))
    conn.commit()
    notify('PSC update',o['desc']+' updated!\n'+o['baseUrl'])
    

def checkUpdates():
    cursor = conn.execute("SELECT ID, baseUrl, username, password, authorization, lastResponse, description from urls")
    rows = cursor.fetchall()
    r=[]
    for row in rows:
        o = {}
        o['ID'] = row[0]
        o['baseUrl'] = row[1]
        o['user'] = row[2]
        o['password'] = row[3]
        o['auth'] = row[4]
        o['lastResponse'] = row[5]
        o['desc'] = row[6]
        r.append(o)
    cursor.close()
    for o in r:
        print(o['desc'])
        checkSingleUpdate(o)

if __name__=='__main__':
    checkUpdates()