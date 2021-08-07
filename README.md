# raftapp
透過 raftexample 改寫，製作一個 raft 協議應用


# feature
1. 基於 mysql 儲存節點狀態
2. 自動擴增節點(加入共識)，也支持手動設定檔復原節點
3. 具備基礎 rest API 功能，`PUT`增加`{"key": "...", "value": "..."}`，`GET`獲取已存的`key`對應值
4. 支持快照(snapshot)介面抽換


# quick start
1. start node1
> ./run s1

2. add & run node2
> ./run s2

3. add & run node3
> ./run s3

4. add key (add '{"key":"jim", "value":"value"}')
> ./run add
<!-- return "ok" -->

5. get key
> ./run get jim
<!-- return "value"> -->
