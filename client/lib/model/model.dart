import 'package:client/api/base.dart';
import 'package:client/common/model.dart';
import 'package:client/model/user.dart';
import 'package:client/tool/date.dart';

/*
* 混合基础结构
* @author wuxubaiyang
* @Time 2022/9/12 19:26
*/
mixin BasePart {
  // id
  late String _id;

  // 更新时间
  late DateTime _updatedAt;

  // 初始化基础结构
  void initBasePart(obj) {
    _id = obj?["id"] ?? "";
    _updatedAt = DateTool.parseDate(obj?["updatedAt"] ?? "") ?? DateTime(0);
  }

  // 获取基础结构map
  Map<String, dynamic> get basePart => {
        "id": id,
        "updatedAt": updatedAt.toString(),
      };

  String get id => _id;

  DateTime get updatedAt => _updatedAt;
}

/*
* 创建者信息
* @author wuxubaiyang
* @Time 2022/10/14 13:52
*/
mixin CreatorPart {
  // 创建者id
  late String _creatorId;

  // 创建者信息
  late UserModel? _creator;

  // 初始化基础结构
  void initCreatorPart(obj) {
    _creatorId = obj?["creatorId"] ?? "";
    if (null != obj?["creator"]) {
      _creator = UserModel.from(obj?["creator"] ?? {});
    }
  }

  // 获取基础结构map
  Map<String, dynamic> get creatorPart => {
        "creatorId": creatorId,
        if (_creator != null) "creator": creator?.to(),
      };

  String get creatorId => _creatorId;

  UserModel? get creator => _creator;
}

/*
* 分页结构
* @author wuxubaiyang
* @Time 2022/9/12 17:22
*/
class PaginationModel<T extends BaseModel> extends BaseModel {
  // 页码
  final int pageIndex;

  // 单页数据量
  final int pageSize;

  // 总数据量
  final int total;

  // 数据集合
  final List<T> data;

  PaginationModel.from(obj, {required OnModelParse<T> itemParse})
      : pageIndex = obj?["pageIndex"] ?? 1,
        pageSize = obj?["pageSize"] ?? 15,
        total = obj?["total"] ?? 0,
        data = (obj?["data"] ?? []).map<T>((e) => itemParse(e)).toList();
}

/*
* 授权信息
* @author wuxubaiyang
* @Time 2022/9/12 19:07
*/
class AuthModel extends BaseModel {
  // 授权key
  final String accessToken;

  // 授权刷新key
  final String refreshToken;

  // 用户信息
  final UserModel user;

  // 检查授权信息是否有效（授权key，刷新key，用户id）
  bool check() =>
      accessToken.isNotEmpty && refreshToken.isNotEmpty && user.id.isNotEmpty;

  AuthModel.from(obj)
      : accessToken = obj?["accessToken"] ?? "",
        refreshToken = obj?["refreshToken"] ?? "",
        user = UserModel.from(obj?["user"] ?? {});

  @override
  Map<String, dynamic> to() => {
        "accessToken": accessToken,
        "refreshToken": refreshToken,
        "user": user.to(),
      };
}
