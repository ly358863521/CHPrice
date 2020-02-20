import json
import pymongo
from pathlib import Path

def citysToProvince():

    city_set = set()
    provinceMap = {}
    provinceMap["不存在"] = []
    provinceMap.values
    count = 0

    with open("province.json","r",encoding="utf-8-sig") as f:
        plist = json.load(f)['provinces']
    with open("CityURL.json","r",encoding="utf-8-sig") as f:
        clist = json.load(f)
    
    for city_letter_list in clist:
        for k in city_letter_list:
            city_set.add(k)
            
    print("收录城市数:",len(city_set))

    for pcitys in plist:
        pName = pcitys["provinceName"]
        if pName in city_set:
            # print(pcitys["provinceName"])
            provinceMap[pName] = [pName]
            count+=1
            city_set.remove(pName)
        else:
            provinceMap[pName] =[]
            for city in pcitys["citys"]:
                cName = list(city.values())[0]
                if cName[-1] =="市":
                    cName = cName[:-1]
                # print(cName)
                if cName in city_set:
                    provinceMap[pName].append(cName)
                    count+=1
                    city_set.remove(cName)
                else:
                    provinceMap["不存在"].append(cName)
                    # print(cName,"not exist")
    
    
    provinceMap["无归属"] = list(city_set)

    # print(city_set)
    print("过滤后城市数:",count)
    print("省份数:",len(provinceMap))
    
    with open("PCity.json","w",encoding="utf-8") as f:
        json.dump(provinceMap,f,ensure_ascii=False)


def saveprovinceMap(db):
    
    with open("PCity.json","r",encoding="utf-8-sig") as f:
        provinceMap = json.load(f)
    province = db["province"]
    for name in provinceMap:
        province.insert_one({"name":name,"citys":provinceMap[name]})
    
    print("省份保存结束")

def saveCityPrice(db):
    Citys = db["Citys"]
    Areas = db["Areas"]
    filePath = Path("../data")
    
    with open("PCity.json","r",encoding="utf-8-sig") as f:
        provinceMap = json.load(f)
    
    for filename in filePath.iterdir():
        with open(filename,"r",encoding="utf-8-sig") as f:
            cityPrice = json.load(f)
            for v in cityPrice["City"].values():

                cityName = v["Name"]
                for province in provinceMap:
                    for i in provinceMap[province]:
                        if i==cityName:
                            v["Province"] = province

                AreaPrice = v["Area"]
                area = []

                if AreaPrice!=[]:
                    for i in range(len(AreaPrice)):
                        AreaPrice[i]["City"] = cityName
                        area.append(AreaPrice[i]["Name"])
                    Areas.insert_many(AreaPrice)

                v["Area"] = area
                Citys.insert_one(v)

    print("城市地区保存结束")

def updateProvincePrice(db):
    province = db["province"]
    citys = db["Citys"]
    for p in province.find():
        price = 0
        tot = 0
        count = 0
        if p["name"]=="不存在" or p["name"] =="无归属":
            continue
        pname = p["name"]
        for city in p["citys"]:
            pricelist = citys.find({"Name":city})[0]["Price"]
            if len(pricelist)>0:
                price = pricelist[-1]
            else:
                price = 0
            if price>0:
                tot+=price
                count+=1
        if count ==0:
            p["price"] = 0
        else:
            p["price"] = tot//count
        province.update_one({"name":pname},{"$set":p})    
    print("省份房价更新完毕")


def savetoMongodb():

    myclient = pymongo.MongoClient("mongodb://localhost:27017/")
    db = myclient["CHPrice"]

    saveprovinceMap(db)
    saveCityPrice(db)
    updateProvincePrice(db)




if __name__ == "__main__":
    # citysToProvince()
    savetoMongodb()
        