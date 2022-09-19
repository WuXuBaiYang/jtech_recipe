import 'dart:io';

import 'package:client/common/api/base.dart';
import 'package:client/common/api/response.dart';
import 'package:client/common/common.dart';
import 'package:client/manage/auth.dart';
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
              onResponse: (res, handler) {
                // 拦截401授权失效
                if (res.statusCode == 401) {
                  /// 判断本地是否存在刷新token，存在则去服务器刷新token
                  /// 不存在token或刷新失败，则路由到登录页重新登录
                  // authManage.refreshToken
                  return handler.reject(
                    DioError(requestOptions: res.requestOptions),
                  );
                }
                return handler.next(res);
              },
              onRequest: (options, handler) {
                var headers = <String, dynamic>{
                  "Platform": Platform.operatingSystem,
                };
                // 已登录则添加token
                if (authManage.authorized) {
                  headers["Authorization"] = "Bearer ${authManage.accessToken}";
                }
                // 添加请求头
                options.headers.addAll(headers);
                handler.next(options);
              },
            ),
          ],
        );

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
