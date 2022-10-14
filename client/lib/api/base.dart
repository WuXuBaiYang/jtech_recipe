import 'package:client/api/auth.dart';
import 'package:client/common/api/base.dart';
import 'package:client/common/api/request.dart';
import 'package:client/common/api/response.dart';
import 'package:client/common/common.dart';
import 'package:client/manage/auth.dart';
import 'package:client/manage/router.dart';
import 'package:dio/dio.dart';

// 数据结构解析回调
typedef OnModelParse<T> = T Function(dynamic obj);

/*
* 本服务api接口基类
* @author wuxubaiyang
* @Time 2022/9/8 16:36
*/
class BaseJAPI extends BaseAPI {
  BaseJAPI()
      : super(
          baseUrl: Common.baseUrl,
          interceptors: [
            InterceptorsWrapper(
              onError: (e, handler) async {
                // 拦截401授权失效
                if (e.response?.statusCode == 401) {
                  // 判断本地是否存在刷新token，存在则去服务器刷新token
                  if (authManage.authorized) {
                    try {
                      await authApi.refreshToken();
                      return handler.next(AuthError.success(e));
                    } catch (_) {}
                  }
                  return handler.next(AuthError.fail(e));
                }
                return handler.next(e);
              },
              onRequest: (options, handler) {
                // 添加请求头
                options.headers.addAll({
                  if (authManage.authorized)
                    "Authorization": "Bearer ${authManage.accessToken}",
                });
                handler.next(options);
              },
            ),
          ],
        );

  @override
  // 重写request方法，实现授权失败业务
  Future<ResponseModel> request(
    String path, {
    RequestModel? request,
    String method = "GET",
    String? cancelKey,
    Options? options,
    OnResponseHandle? responseHandle,
    ProgressCallback? onSendProgress,
    ProgressCallback? onReceiveProgress,
  }) async {
    try {
      return await super.request(path,
          request: request,
          method: method,
          cancelKey: cancelKey,
          options: options,
          responseHandle: responseHandle,
          onSendProgress: onSendProgress,
          onReceiveProgress: onReceiveProgress);
    } on AuthError catch (e) {
      if (e.success) {
        return await super.request(path,
            request: request,
            method: method,
            cancelKey: cancelKey,
            options: options,
            responseHandle: responseHandle,
            onSendProgress: onSendProgress,
            onReceiveProgress: onReceiveProgress);
      }
      // 授权刷新失败，跳转到登录页面
      routerManage.pushNamedAndRemoveUntil(RoutePath.login, untilPath: "");
      rethrow;
    } catch (e) {
      rethrow;
    }
  }

  // 处理报文的状态和数据
  Future<T> handleResponseData<T>(
    Future<ResponseModel> future, {
    OnModelParse<T>? handle,
  }) {
    return future.then<T>((resp) {
      if (resp.success) {
        return handle?.call(resp.data) ?? resp.data;
      }
      throw Exception(resp.message);
    });
  }

  @override
  ResponseModel handleResponse(Response? response) {
    var result = response?.data;
    var code = result?["code"] ?? -1;
    var message = result?["message"] ?? "";
    return ResponseModel(
      code: code,
      message: message,
      success: code == 200,
      data: result?["data"],
    );
  }
}

/*
* 授权异常
* @author wuxubaiyang
* @Time 2022/10/14 16:27
*/
class AuthError extends DioError {
  // 判断授权刷新是否成功
  final bool success;

  AuthError.success(DioError err)
      : success = true,
        super(requestOptions: err.requestOptions);

  AuthError.fail(DioError err)
      : success = false,
        super(requestOptions: err.requestOptions);
}
