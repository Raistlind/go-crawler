package persist

import (
	"testing"
	"GolandProjects/goexercises/crawler/model"
	"gopkg.in/olivere/elastic.v5"
	"context"
	"encoding/json"
)

func TestSave(t *testing.T) {
	expected := model.Profile{
		//Url:  "http://album.zhenai.com/u/108906739",
		//Type: "zhenai",
		//Id:   "108906739",
		//Payload: model.Profile{
		Age:        34,
		Height:     162,
		Weight:     57,
		Income:     "3001-5000元",
		Gender:     "女",
		Name:       "安静的雪",
		Xinzuo:     "牡羊座",
		Occupation: "人事/行政",
		Marriage:   "离异",
		House:      "已购房",
		Hokou:      "山东菏泽",
		Education:  "大学本科",
		Car:        "未购车",
		//},
	}
	id, err := save(expected)

	if err != nil {
		panic(err)
	}

	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	resp, err := client.Get().Index("dating_profile").
		Type("zhenai").Id(id).Do(context.Background())
	if err != nil {
		panic(err)
	}

	t.Logf("%s", resp.Source)

	var actual model.Profile
	err = json.Unmarshal(*resp.Source, &actual)

	if err != nil {
		panic(err)
	}

	if actual != expected {
		t.Errorf("god %v; expected %v", actual, expected)
	}
}
