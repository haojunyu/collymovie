package main

// Movie 电影存储结构体
type Movie struct {
	Idx    string `json:"idx"`    // 排行榜序号
	Title  string `json:"title"`  // 电影名称
	Year   string `json:"year"`   // 电影年份
	Info   string `json:"info"`   // 电影信息
	Rating string `json:"rating"` // 电影排名
	URL    string `json:"url"`    // 电影URL
}
