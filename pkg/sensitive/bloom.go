package sensitive

import (
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/go-ego/gse"
)

type BloomFilter struct {
	filter *bloom.BloomFilter
}

func NewBloomFilter(words []string) *BloomFilter {
	bf := bloom.NewWithEstimates(uint(len(words)), 0.01)
	for _, w := range words {
		bf.AddString(w)
	}
	return &BloomFilter{filter: bf}
}

func (bf *BloomFilter) Contains(word string) bool {
	return bf.filter.TestString(word)
}

type SensitiveFilter struct {
	seg   gse.Segmenter
	bloom *BloomFilter
	words map[string]struct{} // 精确词典，降低布隆误判
}

func NewSensitiveFilter(words []string) *SensitiveFilter {
	var seg gse.Segmenter
	seg.LoadDict()

	// 将敏感词加入分词词典，提升分词准确率
	for _, w := range words {
		seg.AddToken(w, 1000, "sensitive")
	}
	wordMap := make(map[string]struct{}, len(words))
	for _, w := range words {
		wordMap[w] = struct{}{}
	}
	return &SensitiveFilter{
		seg:   seg,
		bloom: NewBloomFilter(words),
		words: wordMap,
	}
}

func (sf *SensitiveFilter) ContainsSensitive(text string) (bool, string) {
	words := sf.seg.Cut(text, true) // 搜索引擎模式分词
	for _, w := range words {
		if sf.bloom.Contains(w) {
			// 二次精确匹配，降低误判
			if _, ok := sf.words[w]; ok {
				return true, w
			}
		}
	}
	return false, ""
}

// 替换文本中的敏感词为指定字符
func (sf *SensitiveFilter) ReplaceSensitive(text string, mask rune) string {
	words := sf.seg.Cut(text, true)
	result := []rune(text)
	for _, w := range words {
		if sf.bloom.Contains(w) {
			if _, ok := sf.words[w]; ok {
				// 替换所有出现的敏感词
				result = []rune(replaceAll(string(result), w, mask))
			}
		}
	}
	return string(result)
}

// 用指定字符替换所有出现的敏感词
func replaceAll(text, word string, mask rune) string {
	runes := []rune(text)
	wordRunes := []rune(word)
	n := len(wordRunes)
	for i := 0; i <= len(runes)-n; i++ {
		if string(runes[i:i+n]) == word {
			for j := range n {
				runes[i+j] = mask
			}
			i += n - 1
		}
	}
	return string(runes)
}