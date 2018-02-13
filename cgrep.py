#!/usr/bin/env python

import json
import base64
import re
import logging
import sys
from urllib import request


def curl_consul(url):
    return request.urlopen(url).read().decode()


def decode_values(json_data):
    decoded_values = dict()
    for j in json_data:
        try:
            decoded_values[j.get('Key')] = base64.b64decode(j.get('Value', "")).decode()
        except Exception as e:
            logging.log(logging.INFO, "Error decoding value of Key: {} - {}".format(j.get('Key'), type(e)))
    return decoded_values


def search_kvs(consul_kvs, pattern):
    found = dict()
    r = re.compile(pattern)
    for key in consul_kvs:        
        try:
            if r.search(consul_kvs[key]) or r.search(key):
                found[key] = consul_kvs[key]
        except:
            logging.log(logging.ERROR, "some error")
    return found


def print_found_kvs(found_kvs):
    for key in found_kvs:
        print("Key: {}, Value: {}".format(key, found_kvs.get(key)))


if __name__ == "__main__":
    data = curl_consul("http://localhost:8500/v1/kv/?recurse=true")
    j = json.loads(data)
    d = decode_values(j)
    
    print_found_kvs(search_kvs(d, sys.argv[1]))

