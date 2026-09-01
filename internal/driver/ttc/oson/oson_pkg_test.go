/*
** Copyright (c) 2026 Oracle and/or its affiliates.
**
** The Universal Permissive License (UPL), Version 1.0
**
** Subject to the condition set forth below, permission is hereby granted to any
** person obtaining a copy of this software, associated documentation and/or data
** (collectively the "Software"), free of charge and under any and all copyright
** rights in the Software, and any and all patent rights owned or freely
** licensable by each licensor hereunder covering either (i) the unmodified
** Software as contributed to or provided by such licensor, or (ii) the Larger
** Works (as defined below), to deal in both
**
** (a) the Software, and
** (b) any piece of software and/or hardware listed in the lrgrwrks.txt file if
** one is included with the Software (each a "Larger Work" to which the Software
** is contributed by such licensors),
**
** without restriction, including without limitation the rights to copy, create
** derivative works of, display, perform, and distribute the Software and make,
** use, sell, offer for sale, import, export, have made, and have sold the
** Software and the Larger Work(s), and to sublicense the foregoing rights on
** either these or other term.
**
** This license is subject to the following condition:
** The above copyright notice and either this complete permission notice or at
** a minimum a reference to the UPL must be included in all copies or
** substantial portions of the Software.
**
** THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
** IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
** FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
** AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
** LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
** OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
** SOFTWARE.
 */

package oson

import (
	"flag"
	"os"
	"strings"
	"testing"
)

// TestCategory selects the test category run by TestCategoryExecutor.
var TestCategory string

func TestMain(m *testing.M) {
	flag.StringVar(&TestCategory, "test.category", "", "testing category, can be unitary, functional, performance, robustness")
	os.Exit(m.Run())
}

var testCases = []struct {
	name       string
	categories string
	exclusive  bool
	f          func(t *testing.T)
}{
	{"TestArrayNode_NestedObjectArrayTraversal", "unitary", false, TestArrayNode_NestedObjectArrayTraversal},
	{"TestArrayNode_RejectsMalformedLayouts", "unitary", false, TestArrayNode_RejectsMalformedLayouts},

	{"TestEncodeStringScalar_UsesExpectedStringOpcodes", "unitary", false, TestEncodeStringScalar_UsesExpectedStringOpcodes},
	{"TestEncodeStringScalar_UsesUB4TreeSegmentSizeWhenTreeExceedsUB2", "unitary", false, TestEncodeStringScalar_UsesUB4TreeSegmentSizeWhenTreeExceedsUB2},
	{"TestEncodeContainers_EncodesNestedObjectAndArray", "unitary", false, TestEncodeContainers_EncodesNestedObjectAndArray},
	{"TestEncodeContainers_EncodesArrayRootWithoutDictionary", "unitary", false, TestEncodeContainers_EncodesArrayRootWithoutDictionary},
	{"TestEncodeContainers_ProducesDeterministicBytesForObjectMaps", "unitary", false, TestEncodeContainers_ProducesDeterministicBytesForObjectMaps},
	{"TestFieldNameSortingOrder", "unitary", false, TestFieldNameSortingOrder},
	{"TestEncodeScalarValues_CoverEssentialScalarOpcodes", "unitary", false, TestEncodeScalarValues_CoverEssentialScalarOpcodes},
	{"TestEncodeScalarValues_SupportsEveryIntegerType", "unitary", false, TestEncodeScalarValues_SupportsEveryIntegerType},
	{"TestEncodeUnsignedInteger_UsesExplicitOracleNumber", "unitary", false, TestEncodeUnsignedInteger_UsesExplicitOracleNumber},
	{"TestEncodeStringNumber_RejectsUB1LengthOverflow", "unitary", false, TestEncodeStringNumber_RejectsUB1LengthOverflow},
	{"TestEncodeContainers_UsesUB2FieldIDs", "unitary", false, TestEncodeContainers_UsesUB2FieldIDs},
	{"TestEncodeContainers_UsesUB4PrimaryDictionaryOffsets", "unitary", false, TestEncodeContainers_UsesUB4PrimaryDictionaryOffsets},
	{"TestEncodeContainers_UsesUB4SecondaryDictionaryOffsets", "unitary", false, TestEncodeContainers_UsesUB4SecondaryDictionaryOffsets},
	{"TestEncodeContainers_UsesAllContainerCountWidths", "unitary", false, TestEncodeContainers_UsesAllContainerCountWidths},
	{"TestEncodeTimestampScalar_UsesTimestampPayload", "unitary", false, TestEncodeTimestampScalar_UsesTimestampPayload},
	{"TestEncodeBinaryFloatAndDoubleScalars_CoverSpecialValues", "unitary", false, TestEncodeBinaryFloatAndDoubleScalars_CoverSpecialValues},
	{"TestEncodeBinaryScalars_CoverLengthBoundaries", "unitary", false, TestEncodeBinaryScalars_CoverLengthBoundaries},
	{"TestEncodeInvalidValues_ReturnOsonEncodingError", "unitary", false, TestEncodeInvalidValues_ReturnOsonEncodingError},
	{"TestEncodeContainers_UsesUB4OffsetsWhenTreeExceedsUB2", "unitary", false, TestEncodeContainers_UsesUB4OffsetsWhenTreeExceedsUB2},
	{"TestEncodeContainers_SupportsLongFieldNames", "unitary", false, TestEncodeContainers_SupportsLongFieldNames},
	{"TestOsonWriteBufferPatchUint_WritesExpectedWidths", "unitary", false, TestOsonWriteBufferPatchUint_WritesExpectedWidths},
	{"TestOsonWriteBufferPatchUint_RejectsInvalidPatch", "unitary", false, TestOsonWriteBufferPatchUint_RejectsInvalidPatch},
	{"TestOsonEncoder_BufferPatchError", "unitary", false, TestOsonEncoder_BufferPatchError},

	{"TestNewNodeAt_RejectsMissingContext", "unitary", false, TestNewNodeAt_RejectsMissingContext},
	{"TestNewNodeAt_ResolvesRedirectChainsAndRejectsCycles", "unitary", false, TestNewNodeAt_ResolvesRedirectChainsAndRejectsCycles},
	{"TestNewNodeAt_RejectsUpdateHeaderOffset", "unitary", false, TestNewNodeAt_RejectsUpdateHeaderOffset},
	{"TestNodeOffsets_CoverAddressWidths", "unitary", false, TestNodeOffsets_CoverAddressWidths},
	{"TestParse_RejectsOutOfRangeRelativeChildOffset", "unitary", false, TestParse_RejectsOutOfRangeRelativeChildOffset},
	{"TestRedirectedNodeOffset_ValidatesEveryMarker", "unitary", false, TestRedirectedNodeOffset_ValidatesEveryMarker},
	{"TestParse_RejectsInvalidUpdateTargets", "unitary", false, TestParse_RejectsInvalidUpdateTargets},
	{"TestParse_RejectsForwardingCycle", "unitary", false, TestParse_RejectsForwardingCycle},
	{"TestOsonDecoderFixtures", "unitary", false, TestOsonDecoderFixtures},
	{"TestOsonDecoder_RejectsNonJSONBinaryFloatText", "unitary", false, TestOsonDecoder_RejectsNonJSONBinaryFloatText},

	{"TestObjectNode_KindReportsObject", "unitary", false, TestObjectNode_KindReportsObject},
	{"TestObjectNode_SimpleObjectTraversal", "unitary", false, TestObjectNode_SimpleObjectTraversal},
	{"TestObjectNode_GetRejectsMalformedChild", "unitary", false, TestObjectNode_GetRejectsMalformedChild},
	{"TestObjectNode_SecondaryDictionaryTraversal", "unitary", false, TestObjectNode_SecondaryDictionaryTraversal},
	{"TestObjectNode_RejectsMalformedLayouts", "unitary", false, TestObjectNode_RejectsMalformedLayouts},
	{"TestObjectNode_SharedOverflowUsesPrimaryTreeOffsets", "unitary", false, TestObjectNode_SharedOverflowUsesPrimaryTreeOffsets},
	{"TestObjectNode_RejectsInvalidDelegateReferences", "unitary", false, TestObjectNode_RejectsInvalidDelegateReferences},
	{"TestReadFieldIDEntriesAt_ReadsAllSupportedWidths", "unitary", false, TestReadFieldIDEntriesAt_ReadsAllSupportedWidths},

	{"TestOsonBuffer_NewBufferInitialState", "unitary", false, TestOsonBuffer_NewBufferInitialState},
	{"TestOsonBuffer_RejectsInvalidInternalCursor", "unitary", false, TestOsonBuffer_RejectsInvalidInternalCursor},
	{"TestOsonBuffer_SetPositionValidatesBounds", "unitary", false, TestOsonBuffer_SetPositionValidatesBounds},
	{"TestOsonBuffer_ReadsSequentialValues", "unitary", false, TestOsonBuffer_ReadsSequentialValues},
	{"TestOsonBuffer_AbsoluteReads", "unitary", false, TestOsonBuffer_AbsoluteReads},
	{"TestOsonBuffer_RejectsInvalidSequentialReads", "unitary", false, TestOsonBuffer_RejectsInvalidSequentialReads},
	{"TestOsonBuffer_RejectsInvalidAbsoluteRanges", "unitary", false, TestOsonBuffer_RejectsInvalidAbsoluteRanges},

	{"TestOsonHeader_SampleOson", "unitary", false, TestOsonHeader_SampleOson},
	{"TestIsOson", "unitary", false, TestIsOson},
	{"TestOsonHeader_ScalarDocumentUsesPostHeaderTreeOffset", "unitary", false, TestOsonHeader_ScalarDocumentUsesPostHeaderTreeOffset},
	{"TestOsonHeader_UpdatedTinyScalarFixture", "unitary", false, TestOsonHeader_UpdatedTinyScalarFixture},
	{"TestOsonHeader_RejectsMalformedUpdateMetadata", "unitary", false, TestOsonHeader_RejectsMalformedUpdateMetadata},
	{"TestOsonHeader_RejectsV1UpdateMetadata", "unitary", false, TestOsonHeader_RejectsV1UpdateMetadata},
	{"TestOsonHeader_RejectsOutOfRangeUpdateMappings", "unitary", false, TestOsonHeader_RejectsOutOfRangeUpdateMappings},
	{"TestOsonHeader_NestedObjectArrayFixture", "unitary", false, TestOsonHeader_NestedObjectArrayFixture},
	{"TestOsonHeader_SecondaryDictionaryFixture", "unitary", false, TestOsonHeader_SecondaryDictionaryFixture},
	{"TestOsonHeader_RejectsMissingInlineLeafFlag", "unitary", false, TestOsonHeader_RejectsMissingInlineLeafFlag},
	{"TestOsonHeader_RejectsTruncatedTreeSegment", "unitary", false, TestOsonHeader_RejectsTruncatedTreeSegment},
	{"TestOsonHeader_RejectsCorruptPrimaryDictionaryOffset", "unitary", false, TestOsonHeader_RejectsCorruptPrimaryDictionaryOffset},
	{"TestOsonHeader_RejectsMalformedFixedHeader", "unitary", false, TestOsonHeader_RejectsMalformedFixedHeader},
	{"TestOsonHeader_RejectsReservedFlags", "unitary", false, TestOsonHeader_RejectsReservedFlags},
	{"TestOsonHeader_ScalarTreeSizeUB4", "unitary", false, TestOsonHeader_ScalarTreeSizeUB4},
	{"TestOsonHeader_ReadHeaderWidthVariants", "unitary", false, TestOsonHeader_ReadHeaderWidthVariants},
	{"TestOsonHeader_MetadataHelpersReflectFlagsAndBounds", "unitary", false, TestOsonHeader_MetadataHelpersReflectFlagsAndBounds},
	{"TestOsonHeader_ForwardingHelpers", "unitary", false, TestOsonHeader_ForwardingHelpers},
	{"TestOsonHeader_AddForwardingAddressValidatesMappings", "unitary", false, TestOsonHeader_AddForwardingAddressValidatesMappings},
	{"TestOsonHeader_RejectsNilBuffer", "unitary", false, TestOsonHeader_RejectsNilBuffer},
	{"TestOsonHeader_RejectsCorruptDictionaryHeaps", "unitary", false, TestOsonHeader_RejectsCorruptDictionaryHeaps},
	{"TestOsonHeader_DictionaryReadersHandleBothOffsetWidths", "unitary", false, TestOsonHeader_DictionaryReadersHandleBothOffsetWidths},
	{"TestOsonHeader_DictionaryReadersRejectTruncation", "unitary", false, TestOsonHeader_DictionaryReadersRejectTruncation},

	{"TestScalarNode_ValueCoversSupportedDecodeUseCases", "unitary", false, TestScalarNode_ValueCoversSupportedDecodeUseCases},
	{"TestScalarNode_NumberAsStringOption", "unitary", false, TestScalarNode_NumberAsStringOption},
	{"TestScalarNode_DefaultOracleNumberAllowsLargePrecisionFloat", "unitary", false, TestScalarNode_DefaultOracleNumberAllowsLargePrecisionFloat},
	{"TestScalarNode_KindAndStringWithOption", "unitary", false, TestScalarNode_KindAndStringWithOption},
	{"TestScalarNode_MalformedScalarPayloads", "unitary", false, TestScalarNode_MalformedScalarPayloads},
	{"TestScalarNode_BinaryFloatSpecialValue", "unitary", false, TestScalarNode_BinaryFloatSpecialValue},
}

func TestCategoryExecutor(t *testing.T) {
	var regularCases, exclusiveCases []struct {
		name       string
		categories string
		exclusive  bool
		f          func(t *testing.T)
	}

	for _, c := range testCases {
		for _, category := range strings.Split(c.categories, ",") {
			if strings.TrimSpace(category) == TestCategory {
				if c.exclusive {
					exclusiveCases = append(exclusiveCases, c)
				} else {
					regularCases = append(regularCases, c)
				}
				break
			}
		}
	}

	if len(regularCases) > 0 {
		t.Run("parallel", func(t *testing.T) {
			t.Parallel()
			for _, c := range regularCases {
				t.Run(c.name, c.f)
			}
		})
	}

	for _, c := range exclusiveCases {
		t.Run(c.name, c.f)
	}
}
