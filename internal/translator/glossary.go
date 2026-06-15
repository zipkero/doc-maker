// 용어집 로드: 사이트 GlossaryPath의 JSON 파일에서 원어→한국어 매핑을 읽는다.
// 포맷: {"원어": "한국어", ...}
//
// 파일이 없거나 비어 있으면 빈 Glossary를 반환하고 오류를 내지 않는다.
// 용어집은 선택적 보조 수단이므로 파일 부재가 번역 중단 사유가 되어서는 안 된다(SPEC §5.6, ANALYSIS D9).
// 단, 파일이 존재하지만 JSON 파싱에 실패하면 오류를 반환한다(잘못된 용어집을 모르고 쓰는 상황 방지).
package translator

import (
	"encoding/json"
	"errors"
	"os"
)

// Glossary는 원어 → 한국어 번역어의 단순 매핑이다(ANALYSIS §5 D9).
// 키: 원어, 값: 한국어 번역어.
type Glossary map[string]string

// LoadGlossary는 path의 JSON 파일을 읽어 Glossary를 반환한다.
//
//   - 파일이 존재하지 않으면 빈 Glossary를 반환한다(오류 없음).
//   - 파일 내용이 비어 있으면 빈 Glossary를 반환한다(오류 없음).
//   - 파일이 존재하지만 JSON 파싱에 실패하면 오류를 반환한다.
//
// 반환된 Glossary가 nil이 되는 경우는 없다.
func LoadGlossary(path string) (Glossary, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		// 파일 없음 → 빈 매핑으로 처리(오류 아님)
		return Glossary{}, nil
	}
	if err != nil {
		return nil, err
	}

	// 빈 파일인 경우 파싱 없이 빈 매핑 반환
	if len(data) == 0 {
		return Glossary{}, nil
	}

	var g Glossary
	if err := json.Unmarshal(data, &g); err != nil {
		return nil, err
	}
	if g == nil {
		return Glossary{}, nil
	}
	return g, nil
}
