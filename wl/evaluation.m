BeginPackage["evaluation`", {
  "MongoDBLink`"
}]


xBegin["`Private`"]


$MonogoDBHost = "csl-224-01.csl.illinois.edu";
$MonogoDBHost = "minsky1-1.csl.illinois.edu";
$MongoDBDatabaseName = "carml";

collections = {
  "evaluation",
  "performance",
  "input_prediction",
  "model_accuracy"
};

conn = OpenConnection[$MonogoDBHost, 27017];

db = GetDatabase[conn, $MongoDBDatabaseName];

evaluationCollection = GetCollection[db, "evaluation"];
modelAccuracyCollection = GetCollection[db, "model_accuracy"];

evaluationCount = CountDocuments[evaluationCollection];

evaluations = Table[
  Association[
    FindDocuments[evaluationCollection, "Offset"->ii, "Limit"->1]
  ]
  ,
  {ii, evaluationCount}
];

$Evaluations = Dataset[evaluations];

accuracyInformation[eval_] :=
  Module[{
    model = Association[eval["model"]],
    modelaccuracyid = eval["modelaccuracyid"],
    modelaccuracy
  },
    If[MissingQ[modelaccuracyid],
      Return[Nothing]
    ];
    modelaccuracy = FindDocuments[modelAccuracyCollection, {"_id" -> modelaccuracyid}, "Limit"->1];
    If[ListQ[modelaccuracy] && Length[modelaccuracy] === 0,
      Return[Nothing]
    ];
    modelaccuracy = Association[First[modelaccuracy]];
    <|
      "Model" -> Lookup[model, "name"],
      "Framework" -> Lookup[Association[Lookup[model, "framework"]], "name"],
      "MachineArchitecture" -> eval["machinearchitecture"],
      "UsingGPU" -> eval["usinggpu"],
      "BatchSize" -> eval["batchsize"],
      "HostName" -> eval["hostname"],
      "Top1" -> modelaccuracy["top1"],
      "Top5" -> modelaccuracy["top5"]
    |>
  ];

$AccuracyInformation = Map[accuracyInformation, evaluations];

CloseConnection[conn];



xEnd[]

EndPackage[]
