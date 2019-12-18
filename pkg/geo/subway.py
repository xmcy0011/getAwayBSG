import requests
from bs4 import BeautifulSoup
import json
"""
目标：爬取中国大陆地铁线路信息
要求：
	①获取相关城市的地铁数量
	②获取每个地铁站的名称
	③写入文档
使用：
python3 subway.py
"""

class Subway(object):
    def __init__(self):
        # 构造url
        self.url = "http://map.amap.com/subway/index.html?&1100"
        # 使用老版本请求头
        self.headers = {
            'user-agent':'Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36'
        }
    # 获取数据
    def get_data(self):
        responses = requests.get(url=self.url, headers=self.headers)
        # 返回str字符串类型
        return responses.text
    # 解析每个城市地铁信息(地铁数量，站点)
    def parse_get_subway(self, ID, city, name):
        # 拼接地铁信息的url
        url = 'http://map.amap.com/service/subway?_1555502190153&srhdata=' + ID + '_drw_' + city + '.json'
        # 获取数据
        response = requests.get(url=url, headers=self.headers)
        # 传递一个参数接收返回的字符串类型
        html = response.text
        # 通过json.loads将json字符串类型转为python数据类型
        result = json.loads(html)
        # 循环遍历数据节点，所有地铁路线
        for node in result['l']:
            # "st"为地铁线的站点
            for start in node['st']:

                # 判断是否含有地铁分线
                # node:"l"里包含所有地铁路线  “la”为分线
                if len(node['la']) > 0:
                    # "ln"为1号线，2号线。。。  “n”为地铁站站名
                    print(name, node['ln'] + '(' + node['la'] + ')', start['n'])

                    with open('subway.json', 'a+', encoding='utf8') as f:
                        f.write(name + ',' + node['ln'] + '(' + node['la'] + ')' + ',' + start['n'] + '\n')

                else:

                    print(name, node['ln'], start['n'])

                    with open('subway.json', 'a+', encoding='utf8') as f:
                        f.write(name + ',' + node['ln'] + ',' + start['n'] + '\n')
    # 解析数据
    def parse_city_data(self, data):
        # 对数据进行编码
        data = data.encode('ISO-8859-1')
        data = data.decode('utf-8')
        soup = BeautifulSoup(data, 'lxml')

        # 获取城市信息
        res1 = soup.find_all(class_="city-list fl")[0]
        res2 = soup.find_all(class_="more-city-list")[0]
        # 遍历a标签
        for temp in res1.find_all('a'):
            # 城市ID值
            ID = temp['id']
            # 城市拼音名
            city_name = temp['cityname']
            # 城市名
            name = temp.get_text()
            self.parse_get_subway(ID, city_name, name)

        for temp in res2.find_all('a'):
            # 城市ID值
            ID = temp['id']
            # 城市拼音名
            city_name = temp['cityname']
            # 城市名
            name = temp.get_text()
            self.parse_get_subway(ID, city_name, name)


    def run(self):
        data =self.get_data()
        self.parse_city_data(data)


if __name__ == '__main__':
    subway = Subway()
    subway.run()