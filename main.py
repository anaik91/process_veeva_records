import os
import sys
import requests
import concurrent.futures

base_url = 'https://sb-covance-etmf-sbx.veevavault.com'
api_version = 'v20.3'
admin_username = 'monisha.desai@'
admin_password = '<password>'
disable = 'Objectlifecyclestateuseraction.study_person__clin.active_state__c.change_state_to_inactive_useraction3__c'
enable = 'Objectlifecyclestateuseraction.study_person__clin.inactive_state__c.change_state_to_active_useraction5__c'

def open_csv(filename):
    with open(filename) as fl:
        data = fl.read()
    return data

def get_session_id(base_url,api_version,username,password):
    url = '{}/api/{}/auth'.format(base_url,api_version)
    headers = {
        'Content-Type': 'application/x-www-form-urlencoded'
    }
    payload='username={}&password={}'.format(username,password)
    r = requests.post(url, headers=headers, data=payload)
    if r.status_code == 200:
        data = r.json()
        return data['sessionId']
    else:
        return None

def change_user_state(base_url,api_version,sessionId,userid,state):
    url ='{}/api/{}/vobjects/study_person__clin/{}/actions/{}'.format(base_url,api_version,userid,state)
    headers = {
    'Content-Type': 'application/x-www-form-urlencoded',
    'Authorization': sessionId
    }
    payload = ''
    r = requests.post(url, headers=headers, data=payload)
    if r.status_code == 200:
        data = r.json()
        if data['responseStatus'] == 'SUCCESS':
            return True
        else:
            return False
    else:
        return False


def process_requests(userids,sessionId,state):
    for user in userids:
        change_user_state(
            base_url,
            api_version,
            sessionId,
            user,
            state
        )

def trigger_concurrently(base_url,api_version,username,password,useridlistchunk,action):
    if action == 'disable':
        state = disable
    else:
        state = enable

    with concurrent.futures.ThreadPoolExecutor(max_workers=5) as executor:
        futures = []
        for userids in useridlistchunk:
            sessionId = get_session_id(base_url,api_version,username,password)
            futures.append(executor.submit(process_requests, userids=userids,sessionId=sessionId,state=state))
        for future in concurrent.futures.as_completed(futures):
            print(future.result())

def main():
    print('Number of CPU Cores: {}'.format(multiprocessing.cpu_count()))
    chunksize = int(os.getenv('CHUNK_SIZE',100))
    print(chunksize)
    userlist = 'users.csv'
    data = open_csv(userlist)
    data = data.splitlines()
    useridlist = [] 
    for user in data:
        userid = user.split(',')[1]
        useridlist.append(userid)
    usercount=len(useridlist)
    print(usercount)
    useridlistchunk = []
    for i in range(0,len(useridlist),chunksize):
        useridlistchunk.append(useridlist[i:i+chunksize])
    
    trigger_concurrently(base_url,api_version,admin_username,admin_password,useridlistchunk,'disable')

    #trigger_concurrently(base_url,api_version,admin_username,admin_password,useridlistchunk,'enable')

if __name__ == '__main__':
    main()
