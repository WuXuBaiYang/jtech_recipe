import 'package:client/api/auth.dart';
import 'package:flutter/material.dart';

/*
* 首页
* @author wuxubaiyang
* @Time 2022/9/8 15:01
*/
class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<StatefulWidget> createState() => _HomePageState();
}

/*
* 首页-状态
* @author wuxubaiyang
* @Time 2022/9/8 15:02
*/
class _HomePageState extends State<HomePage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text("首页"),
      ),
      body: const Center(
        child: Text("这里是首页"),
      ),
      floatingActionButton: FloatingActionButton(
        child: const Icon(Icons.device_hub),
        onPressed: () async {
          // await authApi
          //     .login(userName: "wuxubaiyang", password: "123456")
          //     .then((value) => null);
          // userApi.getSubscribeList(pageIndex: 1, pageSize: 15).then((value) {
          //   print("object");
          // });
        },
      ),
    );
  }
}
