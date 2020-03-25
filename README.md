# goadmin

项目线上聊天 案例  http://goadmin.woaishare.cn
项内容管理系统 案例   http://goadmin.woaishare.cn/admin  测试账号：admin  密码：123456


内容管理系统，可以自定义模型，显示列表字段，分类，验证规则，使用beego开发，也有websocket群里功能,类似thinkphp的onethink

安装运行步骤如下

1，新建一个数据库，取名 goadmin
2，导入数据库文件 goadmin.sql ，文件在goadmin/goadmin.sql
3，修改数据库配置文件 ,文件在goadmin/conf/app.conf

	dbhost = 
	dbport = 3306
	dbname = goadmin
	dbuser = 
	dbpassword = 

4，打包项目  bee pack -be GOOS=linux或者 bee pack -be GOOS=windows
5，项目运行  nohup ./goadmin &


在linux下面的apache配置 httpd.conf如下
<VirtualHost *:80>
    ServerAdmin webmaster@dummy-host.example.com
    ServerName goadmin.woaishare.cn
    ProxyRequests Off
    <Proxy *>
        Order deny,allow
        Allow from all
    </Proxy>
    ProxyPass /ws ws://127.0.0.1:8080/ws/
    ProxyPass / http://127.0.0.1:8080/
    ProxyPassReverse / http://127.0.0.1:8080/
</VirtualHost>
