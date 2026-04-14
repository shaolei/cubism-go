package model

import (
	"encoding/json"
	"testing"
)

func TestModelJsonParsing(t *testing.T) {
	t.Parallel()

	t.Run("complete model3.json", func(t *testing.T) {
		raw := `{
			"Version": 3,
			"FileReferences": {
				"Moc": "model.moc3",
				"Textures": ["texture00.png", "texture01.png"],
				"Physics": "physics.json",
				"Pose": "pose.json",
				"DisplayInfo": "displayInfo.cdi3.json",
				"Expressions": [
					{"Name": "happy", "File": "exp/happy.exp3.json"},
					{"Name": "angry", "File": "exp/angry.exp3.json"}
				],
				"Motions": {
					"Idle": [
						{"File": "motions/idle_01.motion3.json", "FadeInTime": 0.5, "FadeOutTime": 0.5}
					],
					"TapBody": [
						{"File": "motions/tap_body.motion3.json", "Sound": "tap.wav"}
					]
				},
				"UserData": "userdata.json"
			},
			"Groups": [
				{"Target": "Parameter", "Name": "EyeBlink", "Ids": ["ParamEyeLOpen", "ParamEyeROpen"]}
			],
			"HitAreas": [
				{"Id": "HitAreaBody", "Name": "Body"}
			]
		}`

		var m ModelJson
		if err := json.Unmarshal([]byte(raw), &m); err != nil {
			t.Fatalf("failed to parse model json: %v", err)
		}

		if m.Version != 3 {
			t.Errorf("Version = %v, want 3", m.Version)
		}
		if m.FileReferences.Moc != "model.moc3" {
			t.Errorf("Moc = %v, want model.moc3", m.FileReferences.Moc)
		}
		if len(m.FileReferences.Textures) != 2 {
			t.Errorf("Textures len = %v, want 2", len(m.FileReferences.Textures))
		}
		if m.FileReferences.Physics != "physics.json" {
			t.Errorf("Physics = %v, want physics.json", m.FileReferences.Physics)
		}
		if len(m.FileReferences.Expressions) != 2 {
			t.Errorf("Expressions len = %v, want 2", len(m.FileReferences.Expressions))
		}
		if m.FileReferences.Expressions[0].Name != "happy" {
			t.Errorf("Expression[0].Name = %v, want happy", m.FileReferences.Expressions[0].Name)
		}
		if len(m.FileReferences.Motions) != 2 {
			t.Errorf("Motions len = %v, want 2", len(m.FileReferences.Motions))
		}
		if m.FileReferences.Motions["Idle"][0].FadeInTime != 0.5 {
			t.Errorf("Idle[0].FadeInTime = %v, want 0.5", m.FileReferences.Motions["Idle"][0].FadeInTime)
		}
		if m.FileReferences.Motions["TapBody"][0].Sound != "tap.wav" {
			t.Errorf("TapBody[0].Sound = %v, want tap.wav", m.FileReferences.Motions["TapBody"][0].Sound)
		}
		if len(m.Groups) != 1 {
			t.Errorf("Groups len = %v, want 1", len(m.Groups))
		}
		if m.Groups[0].Name != "EyeBlink" {
			t.Errorf("Groups[0].Name = %v, want EyeBlink", m.Groups[0].Name)
		}
		if len(m.Groups[0].Ids) != 2 {
			t.Errorf("Groups[0].Ids len = %v, want 2", len(m.Groups[0].Ids))
		}
		if len(m.HitAreas) != 1 {
			t.Errorf("HitAreas len = %v, want 1", len(m.HitAreas))
		}
		if m.HitAreas[0].Id != "HitAreaBody" {
			t.Errorf("HitAreas[0].Id = %v, want HitAreaBody", m.HitAreas[0].Id)
		}
	})

	t.Run("minimal model3.json", func(t *testing.T) {
		raw := `{
			"Version": 3,
			"FileReferences": {
				"Moc": "model.moc3",
				"Textures": []
			}
		}`

		var m ModelJson
		if err := json.Unmarshal([]byte(raw), &m); err != nil {
			t.Fatalf("failed to parse minimal model json: %v", err)
		}

		if m.Version != 3 {
			t.Errorf("Version = %v, want 3", m.Version)
		}
		if m.FileReferences.Physics != "" {
			t.Errorf("Physics should be empty for minimal json")
		}
		if m.FileReferences.Motions != nil {
			t.Errorf("Motions should be nil for minimal json")
		}
	})

	t.Run("model with motion sync", func(t *testing.T) {
		raw := `{
			"Version": 3,
			"FileReferences": {
				"Moc": "model.moc3",
				"Textures": [],
				"Motions": {
					"Idle": [
						{"File": "idle.motion3.json", "MotionSync": "syncGroup1"}
					]
				}
			}
		}`

		var m ModelJson
		if err := json.Unmarshal([]byte(raw), &m); err != nil {
			t.Fatalf("failed to parse: %v", err)
		}

		if m.FileReferences.Motions["Idle"][0].MotionSync != "syncGroup1" {
			t.Errorf("MotionSync = %v, want syncGroup1", m.FileReferences.Motions["Idle"][0].MotionSync)
		}
	})

	t.Run("missing required fields defaults", func(t *testing.T) {
		raw := `{}`

		var m ModelJson
		if err := json.Unmarshal([]byte(raw), &m); err != nil {
			t.Fatalf("failed to parse empty json: %v", err)
		}

		if m.Version != 0 {
			t.Errorf("Version should default to 0, got %v", m.Version)
		}
		if m.FileReferences.Moc != "" {
			t.Errorf("Moc should default to empty")
		}
	})
}
