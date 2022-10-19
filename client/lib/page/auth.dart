import 'dart:async';
import 'package:client/api/auth.dart';
import 'package:client/common/common.dart';
import 'package:client/common/notifier.dart';
import 'package:client/main.dart';
import 'package:client/manage/router.dart';
import 'package:client/tool/snack.dart';
import 'package:client/tool/tool.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

/*
* 授权页
* @author wuxubaiyang
* @Time 2022/9/8 15:01
*/
class AuthPage extends StatefulWidget {
  const AuthPage({super.key});

  @override
  State<StatefulWidget> createState() => _AuthPageState();
}

/*
* 授权页-状态
* @author wuxubaiyang
* @Time 2022/9/8 15:02
*/
class _AuthPageState extends State<AuthPage> {
  // 表单key
  final formKey = GlobalKey<FormState>();

  // 手机号输入框控制器
  final phoneController = TextEditingController();

  // 验证码输入框控制器
  final smsCodeController = TextEditingController();

  // 手机号校验状态
  final phoneVerifyNotifier = ValueChangeNotifier<bool>(false);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        padding: const EdgeInsets.all(35),
        alignment: Alignment.center,
        child: Form(
          key: formKey,
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              const FlutterLogo(size: 80),
              const SizedBox(height: 65),
              TextFormField(
                autofocus: true,
                controller: phoneController,
                keyboardType: TextInputType.phone,
                textInputAction: TextInputAction.next,
                onChanged: (v) =>
                    phoneVerifyNotifier.setValue(Tool.verifyPhone(v)),
                autovalidateMode: AutovalidateMode.onUserInteraction,
                inputFormatters: [
                  FilteringTextInputFormatter.digitsOnly,
                ],
                decoration: InputDecoration(
                  label: const Text("手机号"),
                  hintText: "000 0000 0000",
                  prefixIcon: const Icon(Icons.phone),
                  suffixIcon: ValueListenableBuilder<bool>(
                    valueListenable: phoneVerifyNotifier,
                    builder: (_, v, __) {
                      return Visibility(
                        visible: v,
                        child: const Icon(Icons.verified_outlined),
                      );
                    },
                  ),
                ),
                validator: (v) {
                  if (v == null || v.isEmpty) {
                    return "手机号不能为空";
                  }
                  if (!Tool.verifyPhone(v)) {
                    return "手机号校验失败";
                  }
                  return null;
                },
              ),
              const SizedBox(height: 28),
              TextFormField(
                maxLength: 4,
                controller: smsCodeController,
                keyboardType: TextInputType.number,
                textInputAction: TextInputAction.done,
                onFieldSubmitted: (v) => _authSaved(),
                inputFormatters: [
                  FilteringTextInputFormatter.digitsOnly,
                ],
                autovalidateMode: AutovalidateMode.onUserInteraction,
                decoration: InputDecoration(
                  counter: const SizedBox(),
                  label: const Text("验证码"),
                  hintText: "0000",
                  prefixIcon: const Icon(Icons.sms),
                  suffixIcon: ValueListenableBuilder<bool>(
                    valueListenable: phoneVerifyNotifier,
                    builder: (_, verifyPhone, __) {
                      return ValueListenableBuilder<int>(
                        valueListenable: smsCountdownNotifier,
                        builder: (_, countdown, __) {
                          var text = countdown > 0
                              ? "验证码已发送(${countdown ~/ 1000})"
                              : "获取验证码";
                          return TextButton(
                            onPressed: verifyPhone && countdown == 0
                                ? () => _sendSMS()
                                : null,
                            child: countdown == -1
                                ? const SizedBox.square(
                                    dimension: 18,
                                    child: CircularProgressIndicator(),
                                  )
                                : Text(text),
                          );
                        },
                      );
                    },
                  ),
                ),
                validator: (v) {
                  if (v == null || v.isEmpty) {
                    return "验证码不能为空";
                  }
                  return null;
                },
              ),
            ],
          ),
        ),
      ),
      floatingActionButton: ValueListenableBuilder<bool>(
        valueListenable: authStateNotifier,
        builder: (_, v, __) {
          return FloatingActionButton(
            onPressed: !v ? () => _authSaved() : null,
            child: v
                ? const SizedBox.square(
                    dimension: 18,
                    child: CircularProgressIndicator(
                      color: Colors.white,
                    ),
                  )
                : const Icon(Icons.done),
          );
        },
      ),
    );
  }

  // 授权请求状态
  final authStateNotifier = ValueChangeNotifier(false);

  // 授权信息保存
  Future<void> _authSaved() async {
    var currState = formKey.currentState;
    if (currState == null || !currState.validate()) return;
    try {
      authStateNotifier.setValue(true);
      var phoneNumber = phoneController.text;
      var code = smsCodeController.text;
      await authApi.auth(phoneNumber: phoneNumber, code: code);
      // 登录成功跳转首页
      routerManage.pushReplacementNamed(RoutePath.home);
    } catch (e) {
      SnackTool.showMessage(context,
          message: debugMode ? e.toString() : "请求授权失败");
      authStateNotifier.setValue(false);
    }
  }

  // 记录短信验证码获取倒计时
  final smsCountdownNotifier = ValueChangeNotifier<int>(0);

  // 发送短信验证码
  Future<void> _sendSMS() async {
    try {
      smsCountdownNotifier.setValue(-1);
      var phoneNumber = phoneController.value.text;
      if (await authApi.sendSMS(phoneNumber: phoneNumber)) {
        _startSmsCountdown();
        if (debugMode) {
          // 开发模式，不会发送验证码，默认使用手机号后四位
          smsCodeController.text =
              phoneNumber.substring(phoneNumber.length - 4);
        }
      }
    } catch (e) {
      SnackTool.showMessage(context, message: "短信验证码发送失败");
      smsCountdownNotifier.setValue(0);
    }
  }

  // 计时器
  Timer? smsCountdownTimer;

  // 短信验证码获取倒计时
  void _startSmsCountdown() {
    var countDown = debugMode ? 1000 * 5 : 1000 * 60;
    smsCountdownNotifier.setValue(countDown);
    smsCountdownTimer = Timer.periodic(
      const Duration(milliseconds: 1000),
      (t) {
        var v = smsCountdownNotifier.value;
        smsCountdownNotifier.setValue(v - 1000);
        if (smsCountdownNotifier.value <= 0) t.cancel();
      },
    );
  }

  @override
  void dispose() {
    super.dispose();
    // 销毁控制器和计时器
    smsCountdownTimer?.cancel();
    phoneController.dispose();
    smsCodeController.dispose();
    phoneVerifyNotifier.dispose();
    smsCountdownNotifier.dispose();
  }
}
