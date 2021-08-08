# deprecate
受限於部署情境，目前無法實作

# 困難原因
1. raft 包會在節點間建立一個 streaming 連線。連線依靠`固定IP`，但是在 k8s 中無法實現
	1. 使用 stateful pod 產生出的 hostname 雖然可於 pod 內識別，但無法於 go app 中被識別
	2. 使用動態給定 ip，會因為節點重起。導致 ip 更動，且更動後後無法更新原有連線 ip
2. 節點群集狀態管理，不適用於頻繁更新服務
	1. 當節點因為更新或重啟因素，致使節點只剩下一個時。此時再進行更新會導致節點完全丟失狀態
	2. 當單一節點通道尚未完全關閉時，可能處理請求尚未完全即被 k8s 消滅。致使處理中請求丟失
3. 使用無狀態式管理，會導致節點 id 無上限遞增。而且也存在請求丟失狀況
	1. 嘗試使用遞增節點 id 方式，做持續新節點部署。放棄舊節點方式，可能導致 id 持續遞增，無法控制
	2. 同 2.2. 問題，可能導致舊節點尚未取消狀態，致使後續叢集發生溝通錯亂

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
