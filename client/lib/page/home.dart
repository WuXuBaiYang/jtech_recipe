import 'package:client/widget/avatar.dart';
import 'package:client/widget/image.dart';
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
        title: const Text('首页'),
      ),
      floatingActionButton: FloatingActionButton(
        child: const Icon(Icons.device_hub),
        onPressed: () async {
          setState(() {});
        },
      ),
    );
  }
}
