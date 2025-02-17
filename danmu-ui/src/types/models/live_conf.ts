//POST /api/live-conf
export interface CreateLiveConfRequest {
  room_display_id: string;
  url: string;
  name: string;
  enable: boolean;
}

//PUT /api/live-conf
export interface UpdateLiveConfRequest {
  id: string;
  room_display_id: string;
  url: string;
  name: string;
  enable: boolean;
}

export interface LiveConf {
  id: string;
  room_display_id: string;
  url: string;
  name: string;
  enable: boolean;
  modified_on: number;
  created_on: number;
  modified_by: string;
  created_by: string;
}

