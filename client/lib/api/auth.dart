import 'package:client/common/api/request.dart';
import 'package:client/manage/auth.dart';
import 'package:client/model/model.dart';
import 'base.dart';

/*
* 授权接口
* @author wuxubaiyang
* @Time 2022/9/12 18:48
*/
class AuthAPI extends BaseJAPI {
  // 获取短信验证码
  Future<bool> sendSMS({required String phoneNumber}) {
    return handleResponseData(
      post("/sms/$phoneNumber"),
    );
  }

  // 请求授权
  Future<AuthModel> auth({
    required String phoneNumber,
    required String code,
  }) {
    return handleResponseData(
      post("/auth",
          requestModel: RequestModel.body(data: {
            "phoneNumber": phoneNumber,
            "code": code,
          })),
      handle: (e) => AuthModel.from(e),
    ).then(
      (v) => Future.wait([
        // 设置授权信息
        authManage.setupAuthInfo(v),
      ]).then((_) => v),
    );
  }

  // 刷新token
  Future<AuthModel> refreshToken() {
    return handleResponseData(
      post("/refreshToken",
          requestModel: RequestModel.create(headers: {
            "RefreshToken": authManage.refreshToken,
          })),
      handle: (e) => AuthModel.from(e),
    ).then(
      (v) => Future.wait([
        // 设置授权信息
        authManage.setupAuthInfo(v),
      ]).then((_) => v),
    );
  }
}

// 单例调用
final authApi = AuthAPI();
