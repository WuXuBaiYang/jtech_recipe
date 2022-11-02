import 'package:client/common/manage.dart';
import 'package:dio/dio.dart';

/*
* 请求撤销管理
* @author wuxubaiyang
* @Time 2022/3/29 14:54
*/
class APICancelManage extends BaseManage {
  static final APICancelManage _instance = APICancelManage._internal();

  factory APICancelManage() => _instance;

  APICancelManage._internal();

  // 缓存接口取消key
  final Map<String, JCancelToken> _cancelKeyMap = {};

  // 生成一个取消授权并返回
  JCancelToken generateToken(String key) {
    if (null != _cancelKeyMap[key]) {
      return _cancelKeyMap[key]!;
    }
    final cancelToken = JCancelToken();
    _cancelKeyMap[key] = cancelToken;
    return cancelToken;
  }

  // 判断请求是否已取消
  bool isCanceled(String key) => _cancelKeyMap[key]?.isCancelled ?? true;

  // 移除并取消请求
  void cancel(String key, {String? reason}) =>
      _cancelKeyMap.remove(key)?.cancel(reason);

  // 取消所有请求
  void cancelAll({String? reason}) => _cancelKeyMap.removeWhere((key, value) {
        value.cancel(reason);
        return true;
      });
}

// 单例调用
final apiCancelManage = APICancelManage();

/*
* 请求撤销token
* @author wuxubaiyang
* @Time 2022/3/29 14:54
*/
class JCancelToken extends CancelToken {}
