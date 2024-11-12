/*
 * Hotel Booking API
 * API for hotel room booking system
 *
 * The version of the OpenAPI document: 1.0.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


package org.openapitools.client.model;

import java.util.Objects;
import java.util.Arrays;
import com.google.gson.TypeAdapter;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.io.IOException;
import java.util.UUID;
import org.threeten.bp.LocalDate;

/**
 * InlineObject
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen", date = "2024-11-08T18:00:57.703777383+03:00[Europe/Moscow]")
public class InlineObject {
  public static final String SERIALIZED_NAME_ROOM_ID = "roomId";
  @SerializedName(SERIALIZED_NAME_ROOM_ID)
  private UUID roomId;

  public static final String SERIALIZED_NAME_GUEST_NAME = "guestName";
  @SerializedName(SERIALIZED_NAME_GUEST_NAME)
  private String guestName;

  public static final String SERIALIZED_NAME_GUEST_EMAIL = "guestEmail";
  @SerializedName(SERIALIZED_NAME_GUEST_EMAIL)
  private String guestEmail;

  public static final String SERIALIZED_NAME_GUEST_PHONE = "guestPhone";
  @SerializedName(SERIALIZED_NAME_GUEST_PHONE)
  private String guestPhone;

  public static final String SERIALIZED_NAME_CHECK_IN = "checkIn";
  @SerializedName(SERIALIZED_NAME_CHECK_IN)
  private LocalDate checkIn;

  public static final String SERIALIZED_NAME_CHECK_OUT = "checkOut";
  @SerializedName(SERIALIZED_NAME_CHECK_OUT)
  private LocalDate checkOut;


  public InlineObject roomId(UUID roomId) {
    
    this.roomId = roomId;
    return this;
  }

   /**
   * Get roomId
   * @return roomId
  **/
  @ApiModelProperty(required = true, value = "")

  public UUID getRoomId() {
    return roomId;
  }


  public void setRoomId(UUID roomId) {
    this.roomId = roomId;
  }


  public InlineObject guestName(String guestName) {
    
    this.guestName = guestName;
    return this;
  }

   /**
   * Get guestName
   * @return guestName
  **/
  @ApiModelProperty(required = true, value = "")

  public String getGuestName() {
    return guestName;
  }


  public void setGuestName(String guestName) {
    this.guestName = guestName;
  }


  public InlineObject guestEmail(String guestEmail) {
    
    this.guestEmail = guestEmail;
    return this;
  }

   /**
   * Get guestEmail
   * @return guestEmail
  **/
  @ApiModelProperty(required = true, value = "")

  public String getGuestEmail() {
    return guestEmail;
  }


  public void setGuestEmail(String guestEmail) {
    this.guestEmail = guestEmail;
  }


  public InlineObject guestPhone(String guestPhone) {
    
    this.guestPhone = guestPhone;
    return this;
  }

   /**
   * Get guestPhone
   * @return guestPhone
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public String getGuestPhone() {
    return guestPhone;
  }


  public void setGuestPhone(String guestPhone) {
    this.guestPhone = guestPhone;
  }


  public InlineObject checkIn(LocalDate checkIn) {
    
    this.checkIn = checkIn;
    return this;
  }

   /**
   * Get checkIn
   * @return checkIn
  **/
  @ApiModelProperty(required = true, value = "")

  public LocalDate getCheckIn() {
    return checkIn;
  }


  public void setCheckIn(LocalDate checkIn) {
    this.checkIn = checkIn;
  }


  public InlineObject checkOut(LocalDate checkOut) {
    
    this.checkOut = checkOut;
    return this;
  }

   /**
   * Get checkOut
   * @return checkOut
  **/
  @ApiModelProperty(required = true, value = "")

  public LocalDate getCheckOut() {
    return checkOut;
  }


  public void setCheckOut(LocalDate checkOut) {
    this.checkOut = checkOut;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    InlineObject inlineObject = (InlineObject) o;
    return Objects.equals(this.roomId, inlineObject.roomId) &&
        Objects.equals(this.guestName, inlineObject.guestName) &&
        Objects.equals(this.guestEmail, inlineObject.guestEmail) &&
        Objects.equals(this.guestPhone, inlineObject.guestPhone) &&
        Objects.equals(this.checkIn, inlineObject.checkIn) &&
        Objects.equals(this.checkOut, inlineObject.checkOut);
  }

  @Override
  public int hashCode() {
    return Objects.hash(roomId, guestName, guestEmail, guestPhone, checkIn, checkOut);
  }


  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class InlineObject {\n");
    sb.append("    roomId: ").append(toIndentedString(roomId)).append("\n");
    sb.append("    guestName: ").append(toIndentedString(guestName)).append("\n");
    sb.append("    guestEmail: ").append(toIndentedString(guestEmail)).append("\n");
    sb.append("    guestPhone: ").append(toIndentedString(guestPhone)).append("\n");
    sb.append("    checkIn: ").append(toIndentedString(checkIn)).append("\n");
    sb.append("    checkOut: ").append(toIndentedString(checkOut)).append("\n");
    sb.append("}");
    return sb.toString();
  }

  /**
   * Convert the given object to string with each line indented by 4 spaces
   * (except the first line).
   */
  private String toIndentedString(Object o) {
    if (o == null) {
      return "null";
    }
    return o.toString().replace("\n", "\n    ");
  }

}

