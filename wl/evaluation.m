BeginPackage["evaluation`", {
  "MongoDBLink`"
}]


xBegin["`Private`"]

$MonogoDBHosts = <|
  "CSL224" -> "csl-224-01.csl.illinois.edu",
  "Minsky" -> "minsky1-1.csl.illinois.edu"
|>;

$MonogoDBHostName = "Minsky";

$MonogoDBHost = $MonogoDBHosts[$MonogoDBHostName];

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
performanceCollection = GetCollection[db, "performance"];
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

accuracyInformation[eval0_] :=
  Module[{
    eval,
    model,
    modelaccuracyid,
    modelaccuracy,
    modelName,
    frameworkName
  },
    eval = Association[eval0];
    model = toAssociation[eval["model"]];
    modelaccuracyid = eval["performanceid"];
    If[MissingQ[modelaccuracyid],
      Return[Nothing]
    ];
    modelaccuracy = FindDocuments[modelAccuracyCollection, {"_id" -> modelaccuracyid}, "Limit"->1];
    If[!ListQ[modelaccuracy] || Length[modelaccuracy] === 0,
      Return[Nothing]
    ];
    modelaccuracy = toAssociation[First[modelaccuracy]];
    modelName = Lookup[model, "name"];
    frameworkName = Lookup[Lookup[model, "framework"], "name"];
    <|
      "ID" -> Lookup[eval, "_id"],
      "Model" -> modelName,
      "Framework" -> frameworkName,
      "FrameworkModel" -> frameworkName <> "::" <> modelName,
      "MachineArchitecture" -> eval["machinearchitecture"],
      "UsingGPU" -> eval["usinggpu"],
      "BatchSize" -> eval["batchsize"],
      "HostName" -> eval["hostname"],
      "Top1" -> modelaccuracy["top1"],
      "Top5" -> modelaccuracy["top5"]
    |>
  ];

$AccuracyInformation = Map[accuracyInformation, evaluations];

debug = Print;

durationInformation[eval0_] :=
  Module[{
    eval,
    model,
    performanceid,
    performance,
    trace,
    predictSpans,
    durations,
    spans,
    modelName,
    frameworkName
  },
    eval = Association[eval0];
    model = toAssociation[eval["model"]];
    performanceid = eval["performanceid"];
    If[MissingQ[performanceid],
      Return[Nothing]
    ];
    performance = FindDocuments[performanceCollection, {"_id" -> performanceid}, "Limit"->1];
    If[!ListQ[performance] || Length[performance] === 0,
      debug["cannot find performanceid = ", performanceid];
      Return[Nothing]
    ];
    debug["found performanceid = ", performanceid];
    trace = First[toAssociation[First[performance]]["trace"]["traces"]];
    If[!AssociationQ[trace],
      debug["performanceid = ", performanceid, " is not an association"];
      Return[Nothing]
    ];
    spans = Flatten[getSpans /@ trace["spans"]];
    predictSpans = Select[spans, #["operationname"] === "Predict" &];
    durations = N[Lookup[predictSpans, "duration"]];
    modelName = Lookup[model, "name"];
    frameworkName = Lookup[Lookup[model, "framework"], "name"];
    <|
      "ID" -> Lookup[eval, "_id"],
      "Model" -> modelName,
      "Framework" -> frameworkName,
      "FrameworkModel" -> frameworkName <> "::" <> modelName,
      "MachineArchitecture" -> eval["machinearchitecture"],
      "UsingGPU" -> eval["usinggpu"],
      "BatchSize" -> eval["batchsize"],
      "HostName" -> eval["hostname"],
      "Durations" -> durations
      (* , "Spans" -> spans *)
    |>
  ];

getSpans[span_] :=
  If[KeyExistsQ[span, "spans"],
    Join[{span}, Catenate[getSpans /@ span["spans"]]],
    {span}
  ];

toAssociation0[e_] := e //.  List[a__Rule] :> Association[a];
toAssociation = GeneralUtilities`ToAssociations;

$DurationInformation = Quiet[Map[durationInformation, evaluations]];

(* CloseConnection[conn]; *)



xEnd[]

EndPackage[]
