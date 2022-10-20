import 'package:client/common/manage.dart';
import 'package:client/model/model.dart';

import 'cache.dart';

/*
* 授权信息管理
* @author wuxubaiyang
* @Time 2022/9/12 21:57
*/
class AuthManage extends BaseManage {
  static final AuthManage _instance = AuthManage._internal();

  factory AuthManage() => _instance;

  AuthManage._internal();

  // 授权信息缓存字段
  static const String _authInfoCacheKey = 'auth_info_cache';

  // 授权信息实体
  AuthModel? _auth;

  @override
  Future<void> init() async {
    final json = cacheManage.getJson(_authInfoCacheKey);
    if (json != null) _auth = AuthModel.from(json);
  }

  // 获取用户信息实体
  AuthModel? get authInfo => _auth;

  // 获取token
  String get accessToken {
    if (!authorized) throw Exception('请在授权之后调用');
    return _auth!.accessToken;
  }

  // 获取刷新token
  String get refreshToken {
    if (!authorized) throw Exception('请在授权之后调用');
    return _auth!.refreshToken;
  }

  // 获取用户id
  String get userId {
    if (!authorized) throw Exception('请在授权之后调用');
    return _auth!.user.id;
  }

  // 判断是否已授权
  bool get authorized => authInfo != null;

  // 设置授权信息
  Future<AuthModel> setupAuthInfo(AuthModel? auth) async {
    if (auth == null) throw Exception('授权信息不能为空');
    if (!auth.check()) throw Exception('授权信息异常');
    if (!await cacheManage.setJsonMap(_authInfoCacheKey, auth.to())) {
      throw Exception('授权信息缓存失败');
    }
    return _auth = auth;
  }

  // 注销授权信息
  Future<void> clearAuthInfo() async {
    await cacheManage.remove(_authInfoCacheKey);
    _auth = null;
  }
}

// 单例调用
final authManage = AuthManage();
