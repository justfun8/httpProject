~~main方法，第一个 sync map 存储用户 ID 和 session 结构体，包括 key 和过期时间。第二个 sync map 存储 key 是赌注 ID ，值是安全的双向有序链表，节点存储用户id，skate，容量保持为 20，加锁保证安全。
post方法，添加节点维持链表。highstake方法输出链表。~~
~~oldmain 方法。第一个map也是用户ID和session，第二个map值是赌注，ID value是map存储用户ID和赌注，Post方法直接存入map ，high stake方法，形成每次动态的堆，弹出前20，再更新map为前20,防止数据太多。~~
之后添加测试， dockerfile
用了AI。:( 

* session部分存储2个sync map，key是id,value是session结构体，和key是sessionKey value是id，可以快速查找id session
*   stake 部分，linklist,stake 是sync Map 存储key 是赌注 ID ，值是安全的双向有序链表
