import urllib.request
from bs4 import BeautifulSoup as bs4

f = open("output/all_course_names.txt", "w")
result = ""
fp = urllib.request.urlopen("https://catalog.utdallas.edu/2020/undergraduate/courses/")
data = fp.read().decode("utf8")
soup = bs4(data, 'html.parser')
tbody = soup.select_one('tbody')
count = 0

for href in tbody.select('tr > td:first-child > a'):
    url = "https://catalog.utdallas.edu"+href['href']
    fp2 = urllib.request.urlopen(url)
    data2 = fp2.read().decode("utf8")
    soup2 = bs4(data2, 'html.parser')
    course_list = soup2.select('.course_title')
    course_name = soup2.select('.course_address')

    print("Getting class: " + href.getText() + " with " + str(len(course_list)) + " courses.")
    idx = 0
    for idx in range(0, min(len(course_list), len(course_name))):
        result += course_name[idx].getText() + "\n"
        count += 1
        idx += 1
    print("Finished getting " + href.getText() + ", running total: " + str(count))
print("Finished scraping. Total course count: " + str(count))

f.write(result)
f.close()