import 'dart:async';

import 'package:client/manage/api_cancel.dart';
import 'package:dio/dio.dart';

import 'request.dart';
import 'response.dart';

// 单次请求的响应回调控制
typedef OnResponseHandle = ResponseModel Function(Response? response);

/*
* 接口方法基类
* @author wuxubaiyang
* @Time 2022/3/29 15:05
*/
abstract class BaseAPI {
  // 网路请求库
  final Dio _dio;

  BaseAPI({
    required baseUrl,
    Map<String, dynamic>? parameters,
    Map<String, dynamic>? headers,
    Duration? connectTimeout,
    Duration? receiveTimeout,
    Duration? sendTimeout,
    int? maxRedirects,
    List<Interceptor> interceptors = const [],
  }) : _dio = Dio(
          BaseOptions(
            baseUrl: baseUrl,
            queryParameters: parameters,
            headers: headers,
            connectTimeout: connectTimeout?.inMilliseconds,
            receiveTimeout: receiveTimeout?.inMilliseconds,
            sendTimeout: sendTimeout?.inMilliseconds,
            maxRedirects: maxRedirects,
          ),
        )..interceptors.addAll([...interceptors]);

  // 附件下载
  Future<ResponseModel> download(
    String path, {
    required String savePath,
    RequestModel? request,
    String method = "GET",
    String? cancelKey,
    Options? options,
    bool deleteOnError = true,
    OnResponseHandle? responseHandle,
    ProgressCallback? onReceiveProgress,
    String lengthHeader = Headers.contentLengthHeader,
  }) {
    // 默认值
    cancelKey ??= path;
    options ??= Options();
    return _handleRequest(
      onRequest: _dio.download(
        path,
        savePath,
        queryParameters: request?.parameters,
        data: request?.data,
        options: options
          ..method ??= method
          ..headers ??= request?.headers,
        cancelToken: apiCancelManage.generateToken(cancelKey),
        onReceiveProgress: onReceiveProgress,
        deleteOnError: deleteOnError,
        lengthHeader: lengthHeader,
      ),
      responseHandle: responseHandle,
    );
  }

  // 基本请求方法
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
    // 默认值
    cancelKey ??= path;
    options ??= Options();
    return _handleRequest(
      onRequest: _dio.request(
        path,
        queryParameters: request?.parameters,
        data: request?.data,
        options: options
          ..method ??= method
          ..headers ??= request?.headers,
        cancelToken: apiCancelManage.generateToken(cancelKey),
        onSendProgress: onSendProgress,
        onReceiveProgress: onReceiveProgress,
      ),
      responseHandle: responseHandle,
    );
  }

  // http-get请求
  Future<ResponseModel> get(
    String path, {
    RequestModel? requestModel,
    String? cancelKey,
    OnResponseHandle? responseHandle,
  }) =>
      request(
        path,
        method: "GET",
        request: requestModel,
        cancelKey: cancelKey,
        responseHandle: responseHandle,
      );

  // http-post请求
  Future<ResponseModel> post(
    String path, {
    RequestModel? requestModel,
    String? cancelKey,
    OnResponseHandle? responseHandle,
  }) =>
      request(
        path,
        method: "POST",
        request: requestModel,
        cancelKey: cancelKey,
        responseHandle: responseHandle,
      );

  // http-put请求
  Future<ResponseModel> put(
    String path, {
    RequestModel? requestModel,
    String? cancelKey,
    OnResponseHandle? responseHandle,
  }) =>
      request(
        path,
        method: "PUT",
        request: requestModel,
        cancelKey: cancelKey,
        responseHandle: responseHandle,
      );

  // http-delete请求
  Future<ResponseModel> delete(
    String path, {
    RequestModel? requestModel,
    String? cancelKey,
    OnResponseHandle? responseHandle,
  }) =>
      request(
        path,
        method: "DELETE",
        request: requestModel,
        cancelKey: cancelKey,
        responseHandle: responseHandle,
      );

  // 处理请求响应
  Future<ResponseModel> _handleRequest({
    required Future<Response> onRequest,
    OnResponseHandle? responseHandle,
  }) async {
    responseHandle ??= handleResponse;
    Response? response;
    try {
      response = await onRequest;
    } on DioError catch (e) {
      response = e.response;
      rethrow;
    }
    return responseHandle(response);
  }

  // 处理请求响应
  ResponseModel handleResponse(Response? response);
}
