import requests
import json
import random

class Populator:
    def __init__(self, events=10, content=5, comments=20, labels=2):
        self.numEvents = events
        self.numContent = content
        self.numComments = comments
        self.numLabels = labels

        self.getWords()
    
    def getWords(self, url="http://svnweb.freebsd.org/csrg/share/dict/words?view=co&content-type=text/plain"):
        response = requests.get(url)
        self.WORDS = response.content.splitlines()
    
    def generate(self):
        for i in range(0, self.numEvents):
            eventName = random.choice(self.WORDS).decode("utf-8")
            url = "http://127.0.0.1:8000/api/v1/{0}"
            resp = requests.get(url.format("event?name={0}".format(eventName)))
            if not resp.status_code == 200:
                raise Exception("Failed to create event")
            eventID = resp.json().get("event").get("event_id")
            for j in range(0, self.numContent):
                title = random.choice(self.WORDS).decode("utf-8")
                resp2 = requests.get(url.format("event/{0}/content?title={1}".format(eventID, title)))
                if not resp2.status_code == 200:
                    raise Exception("Failed to create content")
                contentID = resp2.json().get("event").get("content")[0].get("content_id")
                for k in range(0, self.numComments):
                    comment = " ".join([random.choice(self.WORDS).decode("utf-8") for x in range(0, 16)])
                    body = {"body":comment}
                    resp3 = requests.post(url.format("event/{0}/content/{1}/comment".format(eventID,contentID)), json=body)
                    if not resp3.status_code == 200:
                        raise Exception("Failed to create comment")
                    commentID = resp3.json().get("event").get("content")[0].get("comments")[k].get("comment_id")
                    for l in range(0, self.numLabels):
                        label = random.choice(self.WORDS).decode("utf-8")
                        resp4 = requests.get(url.format("event/{0}/content/{1}/comment/{2}/label?label={3}".format(eventID, contentID, commentID, label)))
                        if not resp4.status_code == 200:
                            raise Exception("Failed to create label")


if __name__ == "__main__":
    pop = Populator()
    pop.generate()
    print("Complete")