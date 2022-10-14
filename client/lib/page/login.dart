import 'package:client/api/auth.dart';
import 'package:client/api/base.dart';
import 'package:client/manage/theme.dart';
import 'package:flutter/material.dart';

/*
* 登录页
* @author wuxubaiyang
* @Time 2022/9/8 15:01
*/
class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  State<StatefulWidget> createState() => _LoginPageState();
}

/*
* 登录页-状态
* @author wuxubaiyang
* @Time 2022/9/8 15:02
*/
class _LoginPageState extends State<LoginPage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: const Center(
        child: Text("这里是登录页"),
      ),
    );
  }
}
