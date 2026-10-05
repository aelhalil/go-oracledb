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
	"fmt"
	"os"
	"testing"

	oracleTest "github.com/oracle/go-oracledb/v26/internal/tests"
)

func TestMain(m *testing.M) {
	if err := oracleTest.InitConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "InitConfig failed: %v\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

var testCases = []oracleTest.CategorizedTestCase{
	// Encoder tests.
	{Name: "TestEncodeStringScalar_UsesExpectedStringOpcodes", Categories: "unitary", Exclusive: false, Fn: TestEncodeStringScalar_UsesExpectedStringOpcodes},
	{Name: "TestEncodeRejectsOSONOver32MiB", Categories: "unitary", Exclusive: false, Fn: TestEncodeRejectsOSONOver32MiB},
	{Name: "TestSortFieldNames_OrdersByHashLengthThenUTF8", Categories: "unitary", Exclusive: false, Fn: TestSortFieldNames_OrdersByHashLengthThenUTF8},
	{Name: "TestEncodeScalarValues_CoverEssentialScalarOpcodes", Categories: "unitary", Exclusive: false, Fn: TestEncodeScalarValues_CoverEssentialScalarOpcodes},
	{Name: "TestEncodeScalarValues_SupportsEveryIntegerType", Categories: "unitary", Exclusive: false, Fn: TestEncodeScalarValues_SupportsEveryIntegerType},
	{Name: "TestEncodeUnsignedInteger_FallsBackToExplicitOracleNumber", Categories: "unitary", Exclusive: false, Fn: TestEncodeUnsignedInteger_FallsBackToExplicitOracleNumber},
	{Name: "TestEncodeContainers_UsesUB4PrimaryDictionaryOffsets", Categories: "unitary", Exclusive: false, Fn: TestEncodeContainers_UsesUB4PrimaryDictionaryOffsets},
	{Name: "TestEncodeContainers_UsesUB4SecondaryDictionaryOffsets", Categories: "unitary", Exclusive: false, Fn: TestEncodeContainers_UsesUB4SecondaryDictionaryOffsets},
	{Name: "TestEncodeArrayChildCount_UsesSmallestWidthAtBoundaries", Categories: "unitary", Exclusive: false, Fn: TestEncodeArrayChildCount_UsesSmallestWidthAtBoundaries},
	{Name: "TestEncodeTimestampScalar_UsesTimestampPayload", Categories: "unitary", Exclusive: false, Fn: TestEncodeTimestampScalar_UsesTimestampPayload},
	{Name: "TestEncodeBinaryFloatAndDoubleScalars_CoverSpecialValues", Categories: "unitary", Exclusive: false, Fn: TestEncodeBinaryFloatAndDoubleScalars_CoverSpecialValues},
	{Name: "TestEncodeBinaryScalars_CoverLengthBoundaries", Categories: "unitary", Exclusive: false, Fn: TestEncodeBinaryScalars_CoverLengthBoundaries},
	{Name: "TestEncodeInvalidValues_ReturnOsonEncodingError", Categories: "unitary", Exclusive: false, Fn: TestEncodeInvalidValues_ReturnOsonEncodingError},
	{Name: "TestEncodeContainers_UsesUB4OffsetsWhenTreeExceedsUB2", Categories: "unitary", Exclusive: false, Fn: TestEncodeContainers_UsesUB4OffsetsWhenTreeExceedsUB2},
	{Name: "TestEncodeContainers_SupportsLongFieldNames", Categories: "unitary", Exclusive: false, Fn: TestEncodeContainers_SupportsLongFieldNames},
	{Name: "TestOsonWriteBufferPatchUint_WritesExpectedWidths", Categories: "unitary", Exclusive: false, Fn: TestOsonWriteBufferPatchUint_WritesExpectedWidths},
	{Name: "TestOsonWriteBufferPatchUint_RejectsInvalidPatch", Categories: "unitary", Exclusive: false, Fn: TestOsonWriteBufferPatchUint_RejectsInvalidPatch},

	// Decoder fixtures and parsing.
	{Name: "TestParse_ParsesAndRendersEveryValidFixture", Categories: "unitary", Exclusive: false, Fn: TestParse_ParsesAndRendersEveryValidFixture},
	{Name: "TestParse_DecodesTimestampTZFixtureToTimeValue", Categories: "unitary", Exclusive: false, Fn: TestParse_DecodesTimestampTZFixtureToTimeValue},
	{Name: "TestOsonDecoder_RejectsNonJSONBinaryFloatText", Categories: "unitary", Exclusive: false, Fn: TestOsonDecoder_RejectsNonJSONBinaryFloatText},

	// Nodes and navigation.
	{Name: "TestNewNodeAt_ResolvesRedirectChainsAndRejectsCycles", Categories: "unitary", Exclusive: false, Fn: TestNewNodeAt_ResolvesRedirectChainsAndRejectsCycles},
	{Name: "TestNodeOffsets_CoverAddressWidths", Categories: "unitary", Exclusive: false, Fn: TestNodeOffsets_CoverAddressWidths},
	{Name: "TestRedirectedNodeOffset_ValidatesEveryMarker", Categories: "unitary", Exclusive: false, Fn: TestRedirectedNodeOffset_ValidatesEveryMarker},
	{Name: "TestParse_RejectsInvalidUpdateTargets", Categories: "unitary", Exclusive: false, Fn: TestParse_RejectsInvalidUpdateTargets},
	{Name: "TestParse_RejectsForwardingCycle", Categories: "unitary", Exclusive: false, Fn: TestParse_RejectsForwardingCycle},

	// Array and object nodes.
	{Name: "TestArrayNode_NestedObjectArrayTraversal", Categories: "unitary", Exclusive: false, Fn: TestArrayNode_NestedObjectArrayTraversal},
	{Name: "TestArrayNode_RejectsMalformedLayouts", Categories: "unitary", Exclusive: false, Fn: TestArrayNode_RejectsMalformedLayouts},
	{Name: "TestObjectNode_KindReportsObject", Categories: "unitary", Exclusive: false, Fn: TestObjectNode_KindReportsObject},
	{Name: "TestObjectNode_SimpleObjectTraversal", Categories: "unitary", Exclusive: false, Fn: TestObjectNode_SimpleObjectTraversal},
	{Name: "TestObjectNode_GetRejectsMalformedChild", Categories: "unitary", Exclusive: false, Fn: TestObjectNode_GetRejectsMalformedChild},
	{Name: "TestObjectNode_SecondaryDictionaryTraversal", Categories: "unitary", Exclusive: false, Fn: TestObjectNode_SecondaryDictionaryTraversal},
	{Name: "TestObjectNode_RejectsMalformedLayouts", Categories: "unitary", Exclusive: false, Fn: TestObjectNode_RejectsMalformedLayouts},
	{Name: "TestObjectNode_SharedOverflowUsesPrimaryTreeOffsets", Categories: "unitary", Exclusive: false, Fn: TestObjectNode_SharedOverflowUsesPrimaryTreeOffsets},
	{Name: "TestObjectNode_RejectsInvalidDelegateReferences", Categories: "unitary", Exclusive: false, Fn: TestObjectNode_RejectsInvalidDelegateReferences},
	{Name: "TestReadFieldIDEntriesAt_ReadsAllSupportedWidths", Categories: "unitary", Exclusive: false, Fn: TestReadFieldIDEntriesAt_ReadsAllSupportedWidths},

	// Scalar node decoding.
	{Name: "TestScalarNode_ValueCoversSupportedDecodeUseCases", Categories: "unitary", Exclusive: false, Fn: TestScalarNode_ValueCoversSupportedDecodeUseCases},
	{Name: "TestScalarNode_ValuePreservesNumberTextInJSONNumberMode", Categories: "unitary", Exclusive: false, Fn: TestScalarNode_ValuePreservesNumberTextInJSONNumberMode},
	{Name: "TestScalarNode_StringQuotesStringValue", Categories: "unitary", Exclusive: false, Fn: TestScalarNode_StringQuotesStringValue},
	{Name: "TestScalarNode_ValueRejectsMalformedPayloads", Categories: "unitary", Exclusive: false, Fn: TestScalarNode_ValueRejectsMalformedPayloads},
	{Name: "TestScalarNode_ValuePreservesBinaryFloatInfinity", Categories: "unitary", Exclusive: false, Fn: TestScalarNode_ValuePreservesBinaryFloatInfinity},
	{Name: "TestScalarNode_IDReadsFullUB1Length", Categories: "unitary", Exclusive: false, Fn: TestScalarNode_IDReadsFullUB1Length},
	{Name: "TestScalarNode_RejectsUnknownOpcode", Categories: "unitary", Exclusive: false, Fn: TestScalarNode_RejectsUnknownOpcode},

	// Buffer bounds and reads.
	{Name: "TestOsonBuffer_NewBufferStartsAtDocumentBeginning", Categories: "unitary", Exclusive: false, Fn: TestOsonBuffer_NewBufferStartsAtDocumentBeginning},
	{Name: "TestOsonBuffer_RejectsSequentialUnderflow", Categories: "unitary", Exclusive: false, Fn: TestOsonBuffer_RejectsSequentialUnderflow},
	{Name: "TestOsonBuffer_SetPositionValidatesBounds", Categories: "unitary", Exclusive: false, Fn: TestOsonBuffer_SetPositionValidatesBounds},
	{Name: "TestOsonBuffer_ReadsSequentialValues", Categories: "unitary", Exclusive: false, Fn: TestOsonBuffer_ReadsSequentialValues},
	{Name: "TestOsonBuffer_ReadsAbsoluteValuesWithoutMovingCursor", Categories: "unitary", Exclusive: false, Fn: TestOsonBuffer_ReadsAbsoluteValuesWithoutMovingCursor},
	{Name: "TestOsonBuffer_RejectsInvalidSequentialReads", Categories: "unitary", Exclusive: false, Fn: TestOsonBuffer_RejectsInvalidSequentialReads},
	{Name: "TestOsonBuffer_RejectsInvalidAbsoluteRanges", Categories: "unitary", Exclusive: false, Fn: TestOsonBuffer_RejectsInvalidAbsoluteRanges},

	// Header parsing and metadata.
	{Name: "TestOsonHeader_ParsesObjectWithPrimaryDictionary", Categories: "unitary", Exclusive: false, Fn: TestOsonHeader_ParsesObjectWithPrimaryDictionary},
	{Name: "TestIsOson", Categories: "unitary", Exclusive: false, Fn: TestIsOson},
	{Name: "TestOsonHeader_ScalarDocumentUsesPostHeaderTreeOffset", Categories: "unitary", Exclusive: false, Fn: TestOsonHeader_ScalarDocumentUsesPostHeaderTreeOffset},
	{Name: "TestOsonHeader_ParsesUpdatedScalarWithOverflowSegment", Categories: "unitary", Exclusive: false, Fn: TestOsonHeader_ParsesUpdatedScalarWithOverflowSegment},
	{Name: "TestOsonHeader_RejectsMalformedUpdateMetadata", Categories: "unitary", Exclusive: false, Fn: TestOsonHeader_RejectsMalformedUpdateMetadata},
	{Name: "TestOsonHeader_RejectsV1UpdateMetadata", Categories: "unitary", Exclusive: false, Fn: TestOsonHeader_RejectsV1UpdateMetadata},
	{Name: "TestOsonHeader_RejectsOutOfRangeUpdateMappings", Categories: "unitary", Exclusive: false, Fn: TestOsonHeader_RejectsOutOfRangeUpdateMappings},
	{Name: "TestOsonHeader_ParsesPrimaryDictionaryWithTinyNodeStats", Categories: "unitary", Exclusive: false, Fn: TestOsonHeader_ParsesPrimaryDictionaryWithTinyNodeStats},
	{Name: "TestOsonHeader_ParsesSecondaryDictionaryWithShortAndLongNames", Categories: "unitary", Exclusive: false, Fn: TestOsonHeader_ParsesSecondaryDictionaryWithShortAndLongNames},
	{Name: "TestOsonHeader_RejectsMissingInlineLeafFlag", Categories: "unitary", Exclusive: false, Fn: TestOsonHeader_RejectsMissingInlineLeafFlag},
	{Name: "TestOsonHeader_RejectsTruncatedTreeSegment", Categories: "unitary", Exclusive: false, Fn: TestOsonHeader_RejectsTruncatedTreeSegment},
	{Name: "TestOsonHeader_RejectsMalformedFixedHeader", Categories: "unitary", Exclusive: false, Fn: TestOsonHeader_RejectsMalformedFixedHeader},
	{Name: "TestOsonHeader_RejectsReservedFlags", Categories: "unitary", Exclusive: false, Fn: TestOsonHeader_RejectsReservedFlags},
	{Name: "TestOsonHeader_ReadsScalarTreeSizeAsUB4", Categories: "unitary", Exclusive: false, Fn: TestOsonHeader_ReadsScalarTreeSizeAsUB4},
	{Name: "TestOsonHeader_ReadHeaderSelectsWidthsFromFlags", Categories: "unitary", Exclusive: false, Fn: TestOsonHeader_ReadHeaderSelectsWidthsFromFlags},
	{Name: "TestOsonHeader_MetadataHelpersReflectFlagsAndBounds", Categories: "unitary", Exclusive: false, Fn: TestOsonHeader_MetadataHelpersReflectFlagsAndBounds},
	{Name: "TestOsonHeader_OffsetResolversAcceptValidAndRejectInvalidOffsets", Categories: "unitary", Exclusive: false, Fn: TestOsonHeader_OffsetResolversAcceptValidAndRejectInvalidOffsets},
	{Name: "TestOsonHeader_AddForwardingAddressValidatesMappings", Categories: "unitary", Exclusive: false, Fn: TestOsonHeader_AddForwardingAddressValidatesMappings},
	{Name: "TestOsonHeader_DictionaryReadersHandleBothOffsetWidths", Categories: "unitary", Exclusive: false, Fn: TestOsonHeader_DictionaryReadersHandleBothOffsetWidths},
	{Name: "TestOsonHeader_DictionaryReadersRejectTruncation", Categories: "unitary", Exclusive: false, Fn: TestOsonHeader_DictionaryReadersRejectTruncation},

	// Malformed input across document structures.
	{Name: "TestOsonHeader_RejectsTruncatedInput", Categories: "unitary", Exclusive: false, Fn: TestOsonHeader_RejectsTruncatedInput},
	{Name: "TestNode_ReadHelpersRejectMalformedInput", Categories: "unitary", Exclusive: false, Fn: TestNode_ReadHelpersRejectMalformedInput},
}

func TestCategoryExecutor(t *testing.T) {
	oracleTest.RunCategoryExecutor(t, oracleTest.TestCategories, testCases)
}
