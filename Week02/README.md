configs下存放配置文件
model下存放数据表映射文件
dao下映射model层文件
service中存放服务模块
biz中根据存放业务文件
routes存放apiGateway相关文件
1、model层封装基础sql相关操作，返回err
2、dao层从数据库获取数据，对于err进行wrap打包
3、service层获取数据，对于err进行withMessage上传
4、biz层进行业务组装，对于err一种方式继续withMessage上传，另一种业务降级
5、routes层日志记录

麻烦老师对工程组织架构进行指点，谢谢。