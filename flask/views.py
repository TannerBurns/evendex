from gevent import monkey
monkey.patch_all()

import json
import requests
import grequests
import sys
import os
import datetime

from flask import Flask, jsonify, redirect, url_for, render_template, request

app=Flask(__name__)
app.secret_key = "supersecretkey"

BASE_IP = "127.0.0.1"
BASE_URL = "http://{}:8000/".format(BASE_IP)

@app.route('/')
def index():
    return redirect(url_for("events"))

@app.route('/events', methods=['GET', 'POST'])
def events():
    eventData = None
    events = None
    eventid = 0

    if request.method == "POST":
        print(request.form)
        if request.form.get("new_event"):
            name = request.form.get("new_event")
            resp = requests.get(BASE_URL+"api/v1/event?name={}".format(name))
            if resp.status_code == 200:
                eventid = int(resp.json().get("event").get("event_id"))
        elif "new_content_" in list(request.form.keys())[0]:
            title = request.form.get(list(request.form.keys())[0])
            eventid = int(list(request.form.keys())[0].split("_")[2])
            resp = requests.get(BASE_URL+"api/v1/event/{0}/content?title={1}".format(eventid, title))
            if resp.status_code == 200:
                eventid = int(resp.json().get("event").get("event_id"))
        elif '{"comment":{"event_id"' in list(request.form.keys())[0]:
            comment = request.form.get(list(request.form.keys())[0])
            dat = json.loads(list(request.form.keys())[0])
            body = {"body": comment}
            url = BASE_URL+"api/v1/event/{0}/content/{1}/comment".format(dat.get("comment").get("event_id"), dat.get("comment").get("content_id"))
            resp = requests.post(url, json=body)
            if resp.status_code == 200:
                eventid = int(resp.json().get("event").get("event_id"))
        elif '{"label":{"event_id"' in list(request.form.keys())[0]:
            label = request.form.get(list(request.form.keys())[0])
            dat = json.loads(list(request.form.keys())[0])
            url = BASE_URL+"api/v1/event/{0}/content/{1}/comment/{2}/label?label={3}".format(dat.get("label").get("event_id"), dat.get("label").get("content_id"), dat.get("label").get("comment_id"), label)
            resp = requests.get(url)
            if resp.status_code == 200:
                eventid = int(dat.get("label").get("event_id"))
        else:
            eventid = int(list(request.form.keys())[0])

    resp = requests.get(BASE_URL+"api/v1/event/{0}".format(eventid))
    if resp.status_code == 200:
        eventData = resp.json()

    resp = requests.get(BASE_URL+"api/v1/events")
    if resp.status_code == 200:
        events = resp.json()

    if eventid == 0 and events.get("count") > 0:
        eventid = events.get("events")[0].get("event_id")

    if eventid > 0:
        resp = requests.get(BASE_URL+"api/v1/event/{0}".format(eventid))
        if resp.status_code == 200:
            eventData = resp.json()

    return render_template("home.html", events=events, event=eventData)

if __name__=='__main__':
    app.run(debug=True)
