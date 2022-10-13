import 'package:client/api/base.dart';
import 'package:client/common/model.dart';
import 'package:client/model/user.dart';
import 'package:client/tool/date.dart';

/*
* 混合基础结构
* @author wuxubaiyang
* @Time 2022/9/12 19:26
*/
mixin BaseInfo {
  // id
  late num _id;

  // 创建事件
  late DateTime _createdAt;

  // 更新时间
  late DateTime _updatedAt;

  // 初始化基础结构
  void initialBaseInfo(obj) {
    _id = obj["id"] ?? 0;
    _createdAt = DateTool.parseDate(obj["createdAt"] ?? "") ?? DateTime(0);
    _updatedAt = DateTool.parseDate(obj["updatedAt"] ?? "") ?? DateTime(0);
  }

  // 获取基础结构map
  Map<String, dynamic> get baseInfoMap => {
        "id": id,
        "createdAt": DateTool.formatDate(DatePattern.fullDateTime, createdAt),
        "updatedAt": DateTool.formatDate(DatePattern.fullDateTime, updatedAt),
      };

  // 判断当前实体是否存在(id==0则认为不存在)
  bool get exist => _id != 0;

  num get id => _id;

  DateTime get createdAt => _createdAt;

  DateTime get updatedAt => _updatedAt;
}

//
mixin CreatorInfo {
  // 创建者id
  late num _creatorId;

  // 创建者信息
  late UserModel? _creator;

  // 初始化基础结构
  void initialCreatorInfo(obj) {
    _creatorId = obj?["creatorId"] ?? 0;
    if (null != obj?["creator"]) {
      _creator = UserModel.from(obj?["creator"] ?? {});
    }
  }

  // 获取基础结构map
  Map<String, dynamic> get creatorInfoMap => {
        "creatorId": creatorId,
        if (_creator != null) "creator": creator?.to(),
      };

  num get creatorId => _creatorId;

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

  // 当前本页数据量
  final int currentSize;

  // 数据集合
  final List<T> data;

  PaginationModel.from(obj, {required OnModelParse<T> parseItem})
      : pageIndex = obj?["pageIndex"] ?? 1,
        pageSize = obj?["pageSize"] ?? 15,
        total = obj?["total"] ?? 0,
        currentSize = obj?["currentSize"] ?? 0,
        data = (obj?["data"] ?? []).map<T>((e) => parseItem(e)).toList();
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

  // 检查授权信息是否有效（授权key，刷新key，用户id，im信息）
  bool check() {
    return accessToken.isNotEmpty &&
        refreshToken.isNotEmpty &&
        user.id != 0 &&
        user.imToken.isNotEmpty &&
        user.imUserId.isNotEmpty;
  }

  @override
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
