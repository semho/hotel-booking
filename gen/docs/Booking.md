

# Booking

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **UUID** |  |  [optional]
**roomId** | **UUID** |  | 
**userId** | **UUID** |  |  [optional]
**guestName** | **String** |  | 
**guestEmail** | **String** |  | 
**guestPhone** | **String** |  |  [optional]
**checkIn** | **LocalDate** |  | 
**checkOut** | **LocalDate** |  | 
**status** | [**StatusEnum**](#StatusEnum) |  |  [optional]
**totalPrice** | **Float** |  |  [optional]
**createdAt** | **OffsetDateTime** |  |  [optional]



## Enum: StatusEnum

Name | Value
---- | -----
PENDING | &quot;PENDING&quot;
CONFIRMED | &quot;CONFIRMED&quot;
CANCELLED | &quot;CANCELLED&quot;
COMPLETED | &quot;COMPLETED&quot;



