// Models generated from schema/models.fbs via FlatBuffers.
// These interfaces match the JSON wire format used by the API.
// If you change a model, edit schema/models.fbs and run: ./schema/generate.sh
//
// The canonical generated FlatBuffers classes live in gen/ts/autograder/.
// These plain interfaces are used for JSON serialization in the frontend.

export interface Lab {
  id: string;
  name: string;
  problem_statement: string;
}

export interface Evaluation {
  type: string;
  name: string;
  actual: string | null;
  status: string;
  points: number;
}

export interface GradeResult {
  evaluations: Evaluation[];
  total_points: number;
}

export interface ApiError {
  error: string;
}
