import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-font-test',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="p-8 bg-white">
      <h1 class="text-2xl font-bold mb-6">Font Test Page</h1>
      
      <div class="space-y-4">
        <div class="border p-4 rounded">
          <h3 class="text-lg font-semibold mb-2">Goudy Serial XBold (Condensed)</h3>
          <p class="label-condensed text-xl">This is Goudy Serial XBold font - ABCDEFGHIJKLMNOPQRSTUVWXYZ</p>
          <p class="label-condensed text-lg">abcdefghijklmnopqrstuvwxyz 0123456789</p>
        </div>

        <div class="border p-4 rounded">
          <h3 class="text-lg font-semibold mb-2">EB Garamond</h3>
          <p class="label-eb-garamond text-xl">This is EB Garamond font - ABCDEFGHIJKLMNOPQRSTUVWXYZ</p>
          <p class="label-eb-garamond text-lg">abcdefghijklmnopqrstuvwxyz 0123456789</p>
        </div>

        <div class="border p-4 rounded">
          <h3 class="text-lg font-semibold mb-2">Label Heading</h3>
          <p class="label-heading text-xl">This is Label Heading font - ABCDEFGHIJKLMNOPQRSTUVWXYZ</p>
          <p class="label-heading text-lg">abcdefghijklmnopqrstuvwxyz 0123456789</p>
        </div>

        <div class="border p-4 rounded">
          <h3 class="text-lg font-semibold mb-2">Label Body</h3>
          <p class="label-body text-xl">This is Label Body font - ABCDEFGHIJKLMNOPQRSTUVWXYZ</p>
          <p class="label-body text-lg">abcdefghijklmnopqrstuvwxyz 0123456789</p>
        </div>

        <div class="border p-4 rounded">
          <h3 class="text-lg font-semibold mb-2">Label Bold</h3>
          <p class="label-bold text-xl">This is Label Bold font - ABCDEFGHIJKLMNOPQRSTUVWXYZ</p>
          <p class="label-bold text-lg">abcdefghijklmnopqrstuvwxyz 0123456789</p>
        </div>

        <div class="border p-4 rounded">
          <h3 class="text-lg font-semibold mb-2">Default Font (Inter)</h3>
          <p class="text-xl">This is default Inter font - ABCDEFGHIJKLMNOPQRSTUVWXYZ</p>
          <p class="text-lg">abcdefghijklmnopqrstuvwxyz 0123456789</p>
        </div>
      </div>

             <div class="mt-8 p-4 bg-gray-100 rounded">
         <h3 class="text-lg font-semibold mb-4">Sample Label Previews</h3>
         
         <div class="flex space-x-4">
           <!-- Channel Label -->
           <div class="w-[384px] h-[480px] border border-black rounded-lg p-2 bg-white">
             <div class="flex justify-between items-center">
               <div class="flex flex-col items-center justify-center text-center font-bold label-condensed border-2 border-black px-2 py-1 w-20 h-20 leading-tight">
                 <div class="text-[10px]">MADE</div>
                 <div class="text-[10px]">IN</div>
                 <div class="text-[10px]">INDIA</div>
               </div>
               <div class="text-right">
                 <div class="text-[12px] font-[700] label-condensed">STEEL AUTHORITY OF INDIA LIMITED</div>
                 <div class="font-bold text-[11px] label-condensed">BHILAI STEEL PLANT</div>
               </div>
             </div>
             <div class="text-center font-bold text-[18px] border-2 border-black mt-2 mb-1 label-heading">
               CHANNEL
             </div>
             <div class="text-center font-bold text-[16px] label-bold mt-2">
               HEAT NO: C103262
             </div>
             <div class="text-center font-bold text-[14px] label-bold mt-1">
               SECTION: CHANNEL 75*40*4.8
             </div>
             <div class="text-center font-bold text-[12px] label-bold mt-1">
               GRADE: IS 2062 EZSOBR
             </div>
           </div>

           <!-- TMT Bar Label -->
           <div class="w-[384px] h-[480px] border border-black rounded-lg p-2 bg-white">
             <div class="flex justify-between items-center">
               <div class="flex flex-col items-center justify-center text-center font-bold label-condensed border-2 border-black px-2 py-1 w-20 h-20 leading-tight">
                 <div class="text-[10px]">MADE</div>
                 <div class="text-[10px]">IN</div>
                 <div class="text-[10px]">INDIA</div>
               </div>
               <div class="text-right">
                 <div class="text-[12px] font-[700] label-condensed">STEEL AUTHORITY OF INDIA LIMITED</div>
                 <div class="font-bold text-[11px] label-condensed">BHILAI STEEL PLANT</div>
               </div>
             </div>
             <div class="text-center font-bold text-[18px] border-2 border-black mt-2 mb-1 label-heading">
               TMT BAR
             </div>
             <div class="text-center font-bold text-[16px] label-bold mt-2">
               HEAT NO: C075400
             </div>
             <div class="text-center font-bold text-[14px] label-bold mt-1">
               SECTION: TMT BAR 25
             </div>
             <div class="text-center font-bold text-[12px] label-bold mt-1">
               GRADE: IS 1786 FE550D
             </div>
           </div>
         </div>
       </div>
    </div>
  `,
  styles: []
})
export class FontTestComponent {} 