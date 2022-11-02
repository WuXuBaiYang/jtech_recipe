import 'package:client/common/model.dart';

/*
* 请求响应实体
* @author wuxubaiyang
* @Time 2022/3/29 14:35
*/
class ResponseModel extends BaseModel {
  // 状态码
  final dynamic code;

  // 描述
  final String message;

  // 返回值
  final dynamic data;

  // 请求是否成功
  final bool success;

  // 构建响应对象
  const ResponseModel({
    required this.code,
    required this.message,
    required this.data,
    required this.success,
  });
}
