# raft-app
透過 raftexample 改寫，製作一個 raft 協議應用

# quick start
1. start node1
> ./app s1

2. add node2
> ./app a2

3. start node2
> ./app s2

4. stop node2
> ./app d2

5. add & start node3
> ./app a3 ; ./app s3

6. stop node3
> ctrl + c

7. recover node3
> ./app r3

8. clean all the nodes
> ./app cls
