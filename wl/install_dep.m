BeginPackage["evaluation`dep`"]

InstallMongoDBLink;

Begin["`Private`"]


InstallMongoDBLink[] :=
  Module[{tmp, dest},
    tmp = URLSave["https://github.com/zbjornson/MongoDBLink/archive/master.zip"];
    dest = FileNameJoin[{$AddOnsDirectory, "Applications"}];
    ExtractArchive[tmp, dest];
    RenameDirectory[FileNameJoin[{dest, "MongoDBLink-master"}], FileNameJoin[{dest, "MongoDBLink"}]];
    DeleteFile[tmp];
    Print["Installed MongoDBLink to " <> dest]
  ];




End[]

EndPackage[]
