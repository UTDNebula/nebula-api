import concurrent.futures
import dotenv
import itertools
import json
import logging
import os
import requests
import sys
import time

from bs4 import BeautifulSoup

# adjust recursion limit for parsing large files with Beautiful Soup
sys.setrecursionlimit(10000) # default is 1000

COURSEBOOK_URL = 'https://coursebook.utdallas.edu'

class CoursebookScraper:

    def __init__(self, ptgsessid: str, congfiguration: dict = None, loggingFile: str = 'scraper.log'):
        # handle object variables
        self.ptgsessid = ptgsessid
        self.configuration = congfiguration
        self.loggingFile = loggingFile

        # initialize logging for scraper
        self.logger = logging.getLogger('Coursebook-Scraper')
        logging.basicConfig(
            format='%(asctime)s - %(levelname)s: %(message)s',
            datefmt='%Y-%m-%d %H:%M:%S',
            filename=self.loggingFile, 
            filemode='w', 
            level=logging.INFO, 
        )

        # document object instantiation
        self.logger.info('Started coursebook scraper')


    def getSearchSpecifications(self) -> dict:
        self.logger.info('Obtaining CourseBook search specifications')
        specifications = {}

        # obtain coursebook index page
        self.logger.info('Getting CourseBook index page')
        try:
            r = requests.get(COURSEBOOK_URL)
        except:
            self.logger.error('Failed to obtain CourseBook index page')
            return {}

        # parse index page for search form
        index = BeautifulSoup(r.content, 'html.parser')
        try:
            form = index.select_one('#guidedsearch')
        except:
            self.logger.error('Failed to parse CourseBook index page for search form')
            return {}

        # parse form fields for field name and options
        try:
            fields = form.find_all('tr')
        except:
            self.logger.error('Failed to parse search form for search fields')
            return {}

        for field in fields:
            try:
                title = field.select_one('th').text
                options = [option['value'] for option in field.find_all('option')]
                specifications[title] = options
            except:
                self.logger.error('Failed to parse a search field from search search form')

        # return obtained specifications
        return specifications


    def downloadCoursebookData(self, queryList: list) -> dict:
        # format query list as a proper query string
        query = 'action=search'
        for specifier in queryList:
            query += '&s[]={0}'.format(specifier)

        # specify CourseBook request header
        coursebookHeader = {
            'Cookie': f'PTGSESSID={self.ptgsessid}',
            'Content-Length': '46',
            'Accept': '*/*',
            'X-Requested_with': 'XMLHttpRequest',
            'User-Agent': 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36',
            'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8',
            'Sec-Gpc': '1',
            'Sec-Fetch-Site': 'same-origin',
            'Sec-Fetch-Mode': 'cors',
            'Sec-Fetch_dest': 'empty',
            'Referer': 'https://coursebook.utdallas.edu/',
        }

        # send CourseBook search query
        # self.logger.info('Sending CourseBook query: {0}'.format(query))
        try:
            response = requests.post(
                'https://coursebook.utdallas.edu/clips/clip-cb11-hat.zog',
                headers = coursebookHeader,
                data = query,
            )
            html = BeautifulSoup(response.content, 'html.parser')
        except:
            self.logger.error('Failed to receive response from query: {0}'.format(query))
            return {}

        # parse response for items returned
        try:
            numItems = html.find_all('b')[0].__repr__().split('(')[1].split(')')[0].split(' ')[0] # this should probably be cleaned up
            if numItems == 'no':
                self.logger.info('0 results obtained from query: {0}'.format(query))
                return {}
        except:
            self.logger.error(f'Failed to parse number of items returned from query: {query}')
            return {}

        # parse response for download link
        try:
            downloadLink = html.find_all('a')[0]['href'].replace('\\','').replace('"','')
            if downloadLink[0:14] == '/reportmonkey/':
                downloadID =  downloadLink.split('/')[-1]
            else:
                raise Exception()
        except:
            self.logger.error(f'Failed to parse download link to {numItems} item(s) from query: {query}')
            return {}

        # create reportmonkey downloadURL from downloadID
        downloadURL = f'https://coursebook.utdallas.edu/reportmonkey/cb11-export/{downloadID}/{downloadID}/json'

        # specify reportmonkey request header
        reportmonkeyHeader = {
            'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9',
            'Accept-Encoding': 'gzip, deflate, br',
            'Accept-Language': 'en-US,en;q=0.9',
            'Connection': 'keep-alive',
            'Cookie': f'PTGSESSID={self.ptgsessid}',
            'Referer': 'https://coursebook.utdallas.edu/',
            'Sec-Fetch-Dest': 'document',
            'Sec-Fetch-Mode': 'navigate',
            'Sec-Fetch-Site': 'same-site',
            'Sec-Fetch-User': '?1',
            'Sec-GPC': '1',
            'Upgrade-Insecure-Requests': '1',
            'User-Agent': 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36',
        }

        # send coursemonkey download request for json
        try:
            download = requests.get(downloadURL, headers=reportmonkeyHeader)
            rawJSON = download.text
        except:
            self.logger.error(f'Failed to download JSON response for query: {query}')
            return {}

        # convert rawJSON to JSON object
        try:
            queryAsJSON = json.loads(rawJSON)
        except:
            self.logger.error(f'Failed to convert JSON to dict for query: {query}')
            return {}

        # return content
        return queryAsJSON


    def downloadFromConfig(self, maxThreads: int = 30) -> dict:
        # convert config dict into queries list
        def dict_product(dictionary):
            keys = dictionary.keys()
            vals = dictionary.values()
            for instance in itertools.product(*vals):
                yield dict(zip(keys, instance))

        queries = [[value for value in query.values() if value != ''] for query in list(dict_product(self.configuration))]

        # start multithreaded download
        data = {'downloads': []}
        with concurrent.futures.ThreadPoolExecutor(max_workers=maxThreads) as executor:
            futures = []
            for query in queries:
                futures.append(executor.submit(self.downloadCoursebookData, query))
            for future in concurrent.futures.as_completed(futures):
                data['downloads'].append(future.result())

        # return data
        self.logger.info('Completed Download From Config')
        return data


if __name__ == '__main__':
    # import environment variables
    dotenv.load_dotenv()

    # instantiate scraper with ptgsessid
    scraper = CoursebookScraper(os.environ['PTGSESSID'])

    # load default config
    with open('configs/{}'.format(os.environ['CONFIG_FILE']), 'r') as file:
        scraper.configuration = json.load(file)

    # download data from coursebook
    result = scraper.downloadFromConfig()

    # write results to data.json
    with open('data/{}.json'.format(int(time.time())),'w') as file:
        json.dump(result, file)