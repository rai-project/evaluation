(* ::Package:: *)

BeginPackage["evaluation`", {
  "JLink`"
}]

xBegin["`Private`"]

ReinstallJava[JVMArguments -> "-Xmx2048m"]

Needs["MongoDBLink`"];

$MonogoDBHosts = <|
  "CSL224" -> "csl-224-01.csl.illinois.edu",
  "Minsky" -> "minsky1-1.csl.illinois.edu",
  "Local" -> "localhost"
|>;

$MonogoDBHostName = "Minsky";

$MonogoDBHost = $MonogoDBHosts[$MonogoDBHostName];

$MongoDBDatabaseName = "carml_frameworktrace_cpu_2";

collections = {
  "evaluation",
  "performance",
  "input_prediction",
  "model_accuracy"
};

getSpans[span_] :=
  If[KeyExistsQ[span, "spans"],
    Join[{span}, Catenate[getSpans /@ span["spans"]]],
    {span}
  ];

toAssociation0[e_] := e //.  List[a__Rule] :> Association[a];
toAssociation1 = GeneralUtilities`ToAssociations;
toAssociation[e_] := Replace[e, List[a__Rule] :> Association[a], {0, Infinity}];


conn = OpenConnection[$MonogoDBHost, 27017];

db = GetDatabase[conn, $MongoDBDatabaseName];

evaluationCollection = GetCollection[db, "evaluation"];
performanceCollection = GetCollection[db, "performance"];
modelAccuracyCollection = GetCollection[db, "model_accuracy"];

evaluationCount = CountDocuments[evaluationCollection];
PrintTemporary["begin evaluations"];
evaluations = Table[
  Association[
    FindDocuments[evaluationCollection, "Offset"->ii, "Limit"->1]
  ]
  ,
  {ii, 0, evaluationCount-1}
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
    If[!AssociationQ[eval],
      Print["unable to set eval ", eval];
      Return[Nothing]
    ];
    model = toAssociation[eval["model"]];
    If[!AssociationQ[eval],
      Print["unable to get model ", eval["model"]];
      Return[Nothing]
    ];
    modelaccuracyid = eval["modelaccuracyid"];
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
      "ModelVersion" -> Lookup[model, "version"],
      "Framework" -> frameworkName,
      "FrameworkModel" -> frameworkName <> "::" <> modelName <> "::" <> Lookup[model, "version"],
      "MachineArchitecture" -> eval["machinearchitecture"],
      "UsingGPU" -> eval["usinggpu"],
      "BatchSize" -> eval["batchsize"],
      "HostName" -> eval["hostname"],
      "Top1" -> modelaccuracy["top1"],
      "Top5" -> modelaccuracy["top5"]
    |>
  ];


Module[{numEvaluations = Length[evaluations]},
DynamicModule[{ii = 0},
  PrintTemporary["Getting accuracy information: ", ProgressIndicator[Dynamic[ii], {1, numEvaluations}]];
  $AccuracyInformation = Map[Function[{eval}, ii++; accuracyInformation[eval]], evaluations]
]];

(* debug = Print; *)

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
    modelName = Lookup[model, "name"];
    (* If[ToLowerCase[modelName] =!= "resnet101",
      Return[Nothing]
    ]; *)
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
    With[{
      traces = toAssociation[First[performance]]["trace"]["traces"]
    },
      trace = If[traces === {} || !ListQ[traces],
        PrintTemporary[performance];
        {},
        First[traces]
      ]
    ];
    If[!AssociationQ[trace],
      debug["performanceid = ", performanceid, " is not an association"];
      Return[Nothing]
    ];
    spans = Flatten[getSpans /@ trace["spans"]];
    predictSpans = Select[spans, ToLowerCase[#["operationname"]] === "predict" &];
    durations = N[Lookup[predictSpans, "duration", {}]];
    frameworkName = Lookup[Lookup[model, "framework"], "name"];
    <|
      "ID" -> Lookup[eval, "_id"],
      "Model" -> modelName,
      "ModelVersion" -> Lookup[model, "version"],
      "Framework" -> frameworkName,
      "FrameworkModel" -> frameworkName <> "::" <> modelName <> "::" <> Lookup[model, "version"],
      "MachineArchitecture" -> eval["machinearchitecture"],
      "UsingGPU" -> eval["usinggpu"],
      "BatchSize" -> eval["batchsize"],
      "HostName" -> eval["hostname"],
      "Durations" -> durations
      (* , "Spans" -> spans *)
    |>
  ];


Module[{numEvaluations = Length[evaluations]},
DynamicModule[{ii = 0},
  PrintTemporary["Getting duration information: ", ProgressIndicator[Dynamic[ii], {1, numEvaluations}]];
  $DurationInformation = Map[Function[{eval}, ii++; durationInformation[eval]], evaluations]
]];

(* CloseConnection[conn]; *)



xEnd[]

EndPackage[]
